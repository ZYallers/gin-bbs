package router

import (
	"net/http"
	"src/abs"
	v100Test "src/controller/v100/test"
	v470Bbs "src/controller/v470/bbs"
	"src/library/expvar"
	"src/library/prometheus"

	"github.com/gin-gonic/gin"
)

var api = &abs.Rest{
	"expvar":  {{Version: "1.0.0+", Method: abs.RestMethod{http.MethodGet: 0}, Handler: expvar.RunningStatsHandler}},
	"metrics": {{Version: "1.0.0+", Method: abs.RestMethod{http.MethodGet: 0}, Handler: prometheus.ServerHandler}},

	"test/index/isok": {{Version: "1.0.0+", Method: abs.RestMethod{http.MethodGet: 1, http.MethodPost: 1}, Handler: func(c *gin.Context) { v100Test.Index(c).CheckOk() }}},
	"test/index/main": {{Version: "1.0.0+", Method: abs.RestMethod{http.MethodGet: 1}, Handler: func(c *gin.Context) { v100Test.Index(c).Main() }}},

	// 470 社区
	"diary/add":              {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.Diary(c).Add() }, Signed: true, Logged: true}},
	"diary/list":             {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.Diary(c).List() }, Signed: true, ParAck: true}},
	"diary/check":            {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.Diary(c).Check() }, Signed: true, ParAck: true}},
	"diary/delete":           {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.Diary(c).Delete() }, Signed: true, Logged: true}},
	"diary/detail":           {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.Diary(c).Detail() }, Signed: true, ParAck: true}},
	"diary/examine":          {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.Diary(c).Examine() }}},
	"circle/icon":            {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.Circle(c).Icon() }, Signed: true, ParAck: true}},
	"circle/detail":          {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.Circle(c).Detail() }, Signed: true, ParAck: true}},
	"circle/examine":         {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.Circle(c).Examine() }}},
	"population/detail":      {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.SpecialPopulation(c).Detail() }, Signed: true}},
	"population/examine":     {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.SpecialPopulation(c).Examine() }}},
	"diary/comment/add":      {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.DiaryComment(c).Add() }, Signed: true, Logged: true}},
	"diary/comment/list":     {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodGet: 1}, Handler: func(c *gin.Context) { v470Bbs.DiaryComment(c).List() }, Signed: true}},
	"diary/comment/template": {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodGet: 1}, Handler: func(c *gin.Context) { v470Bbs.DiaryComment(c).Template() }, Signed: true}},
	"diary/like/save":        {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodPost: 1}, Handler: func(c *gin.Context) { v470Bbs.DiaryLike(c).Save() }, Signed: true, Logged: true}},
	"diary/like/list":        {{Version: "4.7.0+", Method: abs.RestMethod{http.MethodGet: 1}, Handler: func(c *gin.Context) { v470Bbs.DiaryLike(c).List() }, Signed: true, ParAck: true}},
}

func (r *router) GetRestApi() *abs.Rest {
	return api
}
