package tool

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
	"github.com/syyongx/php2go"
	"github.com/techoner/gophp/serialize"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"sort"
	app "src/config"
	"strconv"
	"strings"
	"time"
)

// RandIntn
func RandIntn(max int) int {
	rad := rand.New(rand.NewSource(time.Now().Unix()))
	return rad.Intn(max)
}

// MD5
func MD5(str string) string {
	w := md5.New()
	_, _ = io.WriteString(w, str)
	return fmt.Sprintf("%x", w.Sum(nil))
}

// 字符串首字母转成大写
func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArray := []rune(str)
	if strArray[0] >= 97 && strArray[0] <= 122 {
		strArray[0] -= 32
	}
	return string(strArray)
}

// 字符串首字母转成小写
func StrFirstToLower(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] < 97 || strArry[0] > 122 {
		strArry[0] += 32
	}
	return string(strArry)
}

// PushSimpleMessage
func PushSimpleMessage(msg string, isAtAll bool) {
	if !app.RobotEnable {
		return
	}
	host, _ := os.Hostname()
	text := []string{
		msg + "\n---------------------------",
		"App: " + app.Name,
		"Mode: " + gin.Mode(),
		"Listen: " + *app.HttpServerAddr,
		"HostName: " + host,
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"SystemIP: " + SystemIP(),
		"PublicIP: " + PublicIP(),
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": strings.Join(text, "\n") + "\n",
		},
		"at": map[string]interface{}{
			"isAtAll": isAtAll,
		},
	}
	url := "https://oapi.dingtalk.com/robot/send?access_token=" + app.GracefulRobotToken
	_, _ = NewRequest(url).SetHeaders(map[string]string{"Content-Type": "application/json;charset=utf-8"}).SetPostData(postData).Post()
}

// PushContextMessage
func PushContextMessage(ctx *gin.Context, msg string, reqStr string, stack string, isAtAll bool) {
	if !app.RobotEnable {
		return
	}
	host, _ := os.Hostname()
	text := []string{
		msg + "\n---------------------------",
		"App: " + app.Name,
		"Mode: " + gin.Mode(),
		"Listen: " + *app.HttpServerAddr,
		"HostName: " + host,
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"Url: " + "https://" + ctx.Request.Host + ctx.Request.URL.String(),
		"SystemIP: " + SystemIP(),
		"PublicIP: " + PublicIP(),
		"ClientIP: " + ClientIP(ctx.ClientIP()),
	}
	if reqStr != "" {
		text = append(text, "\nRequest:\n"+strings.ReplaceAll(reqStr, "\n", ""))
	}
	if stack != "" {
		text = append(text, "\nStack:\n"+stack)
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": strings.Join(text, "\n") + "\n",
		},
		"at": map[string]interface{}{
			"isAtAll": isAtAll,
		},
	}
	url := "https://oapi.dingtalk.com/robot/send?access_token=" + app.ErrorRobotToken
	_, _ = NewRequest(url).SetHeaders(map[string]string{"Content-Type": "application/json;charset=utf-8"}).SetPostData(postData).Post()
}

// SystemIP
func SystemIP() string {
	if netInterfaces, err := net.Interfaces(); err == nil {
		for i := 0; i < len(netInterfaces); i++ {
			if (netInterfaces[i].Flags & net.FlagUp) != 0 {
				addrs, _ := netInterfaces[i].Addrs()
				for _, address := range addrs {
					if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
						if ipnet.IP.To4() != nil {
							return ipnet.IP.String()
						}
					}
				}
			}
		}
	}
	return "unknown"
}

// 阻塞式的执行外部shell命令的函数, 等待执行完毕并返回标准输出
func ExecShell(name string, arg ...string) ([]byte, error) {
	// 函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command(name, arg...)

	// 读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为[]byte类型
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run执行c包含的命令，并阻塞直到完成。这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了。
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// GetIPByPconline
func GetIPByPconline(ip string) string {
	var result, url = ip, "http://whois.pconline.com.cn/ipJson.jsp?json=true"
	if ip != "" {
		url += "&ip=" + ip
	}
	resp, err := NewRequest(url).SetTimeOut(1 * time.Second).Get()
	if err != nil || resp.Body == "" {
		return result
	}
	body := mahonia.NewDecoder("GBK").ConvertString(resp.Body)
	if body == "" {
		return result
	}
	info := struct{ IP, Addr string }{}
	if app.Json.Unmarshal([]byte(body), &info) != nil {
		return result
	}
	if info.IP != "" && info.Addr != "" {
		result = fmt.Sprintf("%s %s", info.IP, strings.ReplaceAll(info.Addr, " ", ""))
	}
	return result
}

// PublicIP
func PublicIP() string {
	return GetIPByPconline("")
}

// ClientIP
func ClientIP(ip string) string {
	return GetIPByPconline(ip)
}

// NowMemStats
func NowMemStats() string {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return fmt.Sprintf("Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes) NumGoroutine:%d", ms.Alloc, ms.HeapIdle, ms.HeapReleased, runtime.NumGoroutine())
}

// SafeSendChan
func SafeSendChan(ch chan<- interface{}, value interface{}) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()
	ch <- value
	return false
}

// SafeCloseChan
func SafeCloseChan(ch chan interface{}) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = false
		}
	}()
	close(ch)
	return true
}

// PhpUnserialize
func PhpUnserialize(str string) map[string]interface{} {
	vars := make(map[string]interface{}, 10)
	offset := 0
	strlen := php2go.Strlen(str)
	for offset < strlen {
		if index := strings.Index(php2go.Substr(str, uint(offset), -1), "|"); index < 0 {
			break
		}

		pos := php2go.Strpos(str, "|", offset)
		num := pos - offset

		varname := php2go.Substr(str, uint(offset), num)
		offset += num + 1
		data, _ := serialize.UnMarshal([]byte(php2go.Substr(str, uint(offset), -1)))
		vars[varname] = data

		jsonbyte, _ := serialize.Marshal(data)
		offset += php2go.Strlen(string(jsonbyte))
	}
	return vars
}

// PhpSerialize
func PhpSerialize(vars map[string]interface{}) (str string) {
	for k, v := range vars {
		shal, _ := serialize.Marshal(v)
		str += k + "|" + string(shal)
	}
	return
}

// Nl2br nl2br()
// \n\r, \r\n, \r, \n
func Nl2br(str string, isXhtml bool) string {
	r, n, runes := '\r', '\n', []rune(str)
	var br []byte
	if isXhtml {
		br = []byte("<br />")
	} else {
		br = []byte("<br>")
	}
	skip := false
	length := len(runes)
	var buf bytes.Buffer
	for i, v := range runes {
		if skip {
			skip = false
			continue
		}
		switch v {
		case n, r:
			if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
				buf.Write(br)
				skip = true
				continue
			}
			buf.Write(br)
		default:
			buf.WriteRune(v)
		}
	}
	return buf.String()
}

// InArray in_array()
// haystack supported types: slice, array or map
func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	}
	return false
}

// DrainBody
func DrainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// 经典排序返回a=1&b=1
func SortMapByKey(mp map[string]interface{}) string {
	if len(mp) == 0 {
		return ""
	}
	var newMp = make([]string, 0)
	for k := range mp {
		newMp = append(newMp, k)
	}
	sort.Strings(newMp)
	str := ""
	for _, v := range newMp {
		str += v + "=" + mp[v].(string) + "&"
	}
	return strings.TrimRight(str, "&")
}

// 结构体转为map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// 区间范围内获取随机数
// min 最小值
// max  float64 最大值
// decimalNum  int 返回几位小数点
func RandFloat64(min, max float64, decimalNum int) float64 {
	rand.Seed(time.Now().UnixNano())
	limitFloat64 := rand.Float64()*float64(max-min)*100 + float64(min)*100
	limitStr := strconv.FormatFloat(limitFloat64/100, 'f', decimalNum, 64)
	rankLimit, _ := strconv.ParseFloat(limitStr, 64)
	return rankLimit
}

// OssURL
func OssURL(uri string, bucket ...string) string {
	bkt := app.AliOss.UploadBucket
	if len(bucket) == 1 {
		bkt = bucket[0]
	}
	return `https://` + bkt + `-oss.hxsapp.com/` + uri
}

// 带签名http请求
func HttpRequestWithSign(url string, data map[string]interface{}, params ...interface{}) string {
	var (
		method  = http.MethodPost
		timeout = time.Second * 3
	)
	if len(params) >= 1 {
		method = params[0].(string)
	}
	if len(params) >= 2 {
		timeout = params[1].(time.Duration)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	data["utime"] = timestamp

	hash := md5.New()
	hash.Write([]byte(timestamp + app.TokenKey))
	md5str := hex.EncodeToString(hash.Sum(nil))
	data["sign"] = base64.StdEncoding.EncodeToString([]byte(md5str))

	headers := map[string]string{
		"Connection":   "close",
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36",
	}
	req := NewRequest(url).SetMethod(method).SetHeaders(headers).SetTimeOut(timeout)

	switch method {
	case http.MethodGet:
		queries := make(map[string]string, len(data))
		for k, v := range data {
			queries[k] = fmt.Sprintf("%v", v)
		}
		req.SetQueries(queries)
	default:
		req.SetPostData(data)
	}
	if resp, err := req.Send(); err == nil {
		return resp.Body
	}
	return ""
}

// mysql time.time  转换　字符串时间
func ParseTime(timeStr, formatStr string) string {
	p, _ := time.Parse(formatStr, timeStr)
	return p.Format("2006-01-02 15:04:05")
}

// 左边补零
func FixPad(id string) string {
	return fmt.Sprintf("%08s", id)
}

// 左边去零
func FixTrim(id string) string {
	return strings.TrimLeft(id, "0")
}

// 并集
func UnionSlice(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// 删前缀
func DealPreOss(url string, dType, addType int) string {
	url = strings.ToLower(url)
	preStr := "https://hxsupload-oss.hxsapp.com/"
	if dType == 2 {
		preStr = "http://hxsapp-user-media-out-oss.hxsapp.com/"
	}
	if addType == 1 {
		if strings.Contains(url, "http") {
			return url
		}
		return preStr + url
	}
	return strings.Replace(url, preStr, "", 1)
}
