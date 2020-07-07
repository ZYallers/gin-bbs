package app

import "os"

// AliOss 配置
type ossEnv struct {
	AccessKeyId, AccessKeySecret, EndPoint, UploadBucket string
}

var AliOss = ossEnv{
	AccessKeyId:     os.Getenv("alioss_accesskeyid"),
	AccessKeySecret: os.Getenv("alioss_accesskeysecret"),
	EndPoint:        os.Getenv("alioss_endpoint"),
	UploadBucket:    os.Getenv("alioss_uploadbucket"),
}
