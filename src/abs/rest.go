package abs

import "github.com/gin-gonic/gin"

type Rest map[string][]RestHandler

type RestMethod map[string]byte

type RestHandler struct {
	Version                string
	Method                 RestMethod
	Handler                gin.HandlerFunc
	Signed, Logged, ParAck bool
}
