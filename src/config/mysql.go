package app

import (
	"os"
)

type dialect struct {
	User, Pwd, Host, Port, Db string
}

var (
	mysqlPort = os.Getenv("mysql_port")
	// DB 数据库连接配置
	DB        = map[string]dialect{
		"et": {
			Host: os.Getenv("mysql_host"),
			User: os.Getenv("mysql_username"),
			Pwd:  os.Getenv("mysql_password"),
			Db:   os.Getenv("mysql_database"),
			Port: mysqlPort,
		},
		"bbs": {
			Host: os.Getenv("mysql_bbs_host"),
			User: os.Getenv("mysql_bbs_username"),
			Pwd:  os.Getenv("mysql_bbs_password"),
			Db:   os.Getenv("mysql_bbs_database"),
			Port: mysqlPort,
		},
	}
)
