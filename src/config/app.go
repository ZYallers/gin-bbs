package app

import (
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"go.uber.org/zap"
	"time"
)

// 项目配置
const (
	Name                  = "gin-bbs"
	Version               = "4.7.0"
	HttpServerDefaultAddr = "0.0.0.0:9120"
	HttpServerIdleTimeout = 30 * time.Second
	LogDir                = "/apps/logs/go/gin-bbs"
	ErrorRobotToken       = ""
	GracefulRobotToken    = ""

	ReqStrKey = "gin-gonic/gin/reqstr"
	TokenKey  = ""
)

var (
	HttpServerAddr *string
	Engine         *gin.Engine
	Logger         *zap.Logger
	DebugStack     bool
	RobotEnable    bool
	Json           = jsoniter.ConfigCompatibleWithStandardLibrary
)
