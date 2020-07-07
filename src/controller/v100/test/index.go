package v110

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "src/abs"
    app "src/config"
    "src/library/tool"
)

type index struct {
    abs.Controller
}

func Index(ctx *gin.Context) *index {
    i := &index{}
    i.Ctx = ctx
    return i
}

func (i *index) CheckOk() {
    i.Json(http.StatusOK, gin.H{
        "code": http.StatusOK,
        "msg":  "ok",
        "data": gin.H{
            "mode":        gin.Mode(),
            "public_ip":   tool.PublicIP(),
            "system_ip":   tool.SystemIP(),
            "client_ip":   tool.ClientIP(i.Ctx.ClientIP()),
            "request_url": i.Ctx.Request.URL.String(),
        },
    })
}

func (i *index) Main() {
    app.Engine.LoadHTMLFiles("src/view/v100/test/index/main.html")
    i.Ctx.HTML(http.StatusOK, "main.html", gin.H{
        "title": i.Ctx.Request.URL.String(),
    })
}
