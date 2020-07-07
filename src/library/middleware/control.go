package middleware

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/syyongx/php2go"
	"net/http"
	"regexp"
	"src/abs"
	app "src/config"
	"src/library/tool"
	"strconv"
	"strings"
	"time"
)

var (
	session = &abs.Redis{}
)

func Dispatch(ege *gin.Engine, api *abs.Rest) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if regexp.MustCompile(`^\/v[0-9]{3,6}\/.*$`).MatchString(ctx.Request.URL.Path) {
			var handler *abs.RestHandler
			if handler, _ = versionCompare(strings.Join(strings.Split(ctx.Request.URL.Path, `/`)[2:], `/`), ctx, *api); handler == nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "msg": "bad request exception"})
				return
			}
			// 签名验证
			if handler.Signed && !signCheck(ctx) {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "msg": "signature error"})
				return
			}
			// 登录验证
			if handler.Logged {
				if !loginCheck(ctx) {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusUnauthorized, "msg": "please log in and operate again"})
					return
				}
			} else {
				if handler.ParAck {
					ackParsing(ctx)
				}
			}
			ctx.Next()
			go regenSessionData(ctx.Copy())
		} else {
			// 版本验证
			if handler, version := versionCompare(ctx.Request.URL.Path[1:], ctx, *api); handler != nil {
				ctx.Request.URL.Path = "/v" + strings.Join(strings.Split(version, "."), "") + ctx.Request.URL.Path
				ege.HandleContext(ctx)
				ctx.Abort()
			} else {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "msg": "page not found"})
			}
		}
	}
}

// queryPostForm
func queryPostForm(ctx *gin.Context, keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	if val, ok := ctx.GetQuery(keys[0]); ok {
		return val
	}
	if val, ok := ctx.GetPostForm(keys[0]); ok {
		return val
	}
	if len(keys) == 2 {
		return keys[1]
	}
	return ""
}

func versionCompare(restful string, ctx *gin.Context, api abs.Rest) (*abs.RestHandler, string) {
	if handlers, ok := api[restful]; ok {
		version := queryPostForm(ctx, "app_version", app.Version)
		for _, handler := range handlers {
			if _, ok := handler.Method[ctx.Request.Method]; ok {
				if ilene := len(handler.Version); handler.Version[ilene-1:] == "+" {
					vs := handler.Version[0 : ilene-1]
					if php2go.VersionCompare(version, vs, ">=") {
						return &handler, vs
					}
				} else {
					if php2go.VersionCompare(version, handler.Version, "=") {
						return &handler, version
					}
				}
			}
		}
	}
	return nil, ""
}

func signCheck(ctx *gin.Context) (pass bool) {
	sign := queryPostForm(ctx, "sign")
	if sign == "" {
		return
	}
	timestampStr := queryPostForm(ctx, "utime")
	if timestampStr == "" {
		return
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 0)
	if err != nil {
		return
	}
	if time.Now().Unix()-timestamp > 3600 {
		return
	}
	sign = strings.Trim(sign, " ")
	h := md5.New()
	h.Write([]byte(strconv.FormatInt(timestamp, 10) + app.TokenKey))
	md5str := hex.EncodeToString(h.Sum(nil))
	input := []byte(md5str)
	if sign == base64.StdEncoding.EncodeToString(input) {
		pass = true
	}
	return
}

func loginCheck(ctx *gin.Context) (pass bool) {
	if sessToken := queryPostForm(ctx, "sess_token"); sessToken != "" {
		ctx.Set(app.SessionTokenKey, sessToken)
		if vars := getSessionData(sessToken); len(vars) > 0 {
			ctx.Set(app.SessionDataKey, vars)
			if userInfo, ok := vars["userinfo"].(map[string]interface{}); ok {
				if userId, ok := userInfo["userid"].(string); ok && userId != "" {
					pass = true
					ctx.Set(app.SessionLoggedUidKey, userId)
				}
			}
		}
	}
	return
}

func ackParsing(ctx *gin.Context) {
	if sessToken := queryPostForm(ctx, "sess_token"); sessToken != "" {
		ctx.Set(app.SessionTokenKey, sessToken)
		if vars := getSessionData(sessToken); len(vars) > 0 {
			ctx.Set(app.SessionDataKey, vars)
			if userInfo, ok := vars["userinfo"].(map[string]interface{}); ok {
				if userId, ok := userInfo["userid"].(string); ok && userId != "" {
					ctx.Set(app.SessionLoggedUidKey, userId)
				}
			}
		}
	}
}

func getSessionData(sessToken string) (vars map[string]interface{}) {
	if str, _ := session.GetSession().Get(app.Redis.Key.Session.StringSessToken + sessToken).Result(); str != "" {
		vars = tool.PhpUnserialize(str)
	}
	return
}

func regenSessionData(ctx *gin.Context) {
	var (
		sessToken string
		vars      map[string]interface{}
	)

	if value, ok := ctx.Get(app.SessionTokenKey); !ok {
		return
	} else {
		sessToken = value.(string)
	}

	if value, ok := ctx.Get(app.SessionDataKey); !ok {
		return
	} else {
		vars = value.(map[string]interface{})
	}

	nowTime := time.Now()
	if lastRegen, ok := vars["__ci_last_regenerate"].(int); ok {
		if nowTime.After(time.Unix(int64(lastRegen), 0).Add(app.SessionUpdateDuration)) {
			vars["__ci_last_regenerate"] = nowTime.Unix()
			newCiVars := make(map[string]interface{}, 10)
			if ciVars, ok := vars["__ci_vars"].(map[string]interface{}); ok {
				for k := range ciVars {
					newCiVars[k] = nowTime.Unix() + app.SessionExpiration
				}
				vars["__ci_vars"] = newCiVars
			}
			session.GetSession().Set(app.Redis.Key.Session.StringSessToken+sessToken, tool.PhpSerialize(vars), app.SessionExpiration*time.Second)
		}
	}
}
