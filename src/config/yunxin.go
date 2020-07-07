package app

import "os"

type yunXinEnv struct {
	AppKey, AppSecret, SignKey string
}

var YunXin = yunXinEnv{
	AppKey:    os.Getenv("yunxin_appkey"),
	AppSecret: os.Getenv("yunxin_appsecret"),
	SignKey:   os.Getenv("yunxin_signkey"),
}
