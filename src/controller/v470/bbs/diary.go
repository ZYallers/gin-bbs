package v470

import (
	"net/http"
	"src/abs"
	app "src/config"
	"src/library/tool"
	"src/service"
	"src/table"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type diary struct {
	abs.Controller
}

func Diary(c *gin.Context) *diary {
	d := &diary{}
	d.Ctx = c
	return d
}

// 新增动态
func (d *diary) Add() {
	userId := d.GetLoggedUserId()
	content := d.GetQueryByMethod("content", "")
	subData := d.GetQueryByMethod("sub_data", "")
	subType, _ := strconv.Atoi(d.GetQueryByMethod("sub_type", "1"))
	dataType, _ := strconv.Atoi(d.GetQueryByMethod("data_type", "1"))
	dataId, _ := strconv.Atoi(d.GetQueryByMethod("data_id", "0"))
	if dataId <= 0 {
		d.Json(http.StatusNotImplemented, "缺少必要参数data_id")
		return
	}
	if subData != "" {
		dataSlice := []string{}
		for _, v := range strings.Split(subData, ",") {
			if subType == 1 {
				dataSlice = append(dataSlice, tool.DealPreOss(v, 1, 0))
			} else {
				dataSlice = append(dataSlice, tool.DealPreOss(v, 2, 0))
			}
		}
		dataJson, _ := app.Json.Marshal(dataSlice)
		subData = string(dataJson)
	}
	diary := table.BbsDiary{
		UserId:     userId,
		Content:    content,
		SubData:    subData,
		SubType:    subType,
		DataId:     dataId,
		DataType:   dataType,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	if service.NewBbsDiary().AddDiary(diary) > 0 {
		d.Json(http.StatusOK, "发布成功")
	} else {
		d.Json(http.StatusNotImplemented, "发布失败")
	}
}

// 审核动态 更新缓存
func (d *diary) Examine() {
	if ids := d.GetQueryByMethod("ids", ""); ids != "" {
		service.NewBbsDiary().ExamineDiary(strings.Split(ids, ","))
	}
	d.Json(http.StatusOK, "更新缓存成功")
}

// 动态列表
func (d *diary) List() {
	userId := d.GetLoggedUserId()
	dataId, _ := strconv.Atoi(d.GetQueryByMethod("data_id", "0"))
	pageDepend, _ := strconv.Atoi(d.GetQueryByMethod("page_depend", "1"))
	preType, _ := strconv.Atoi(d.GetQueryByMethod("pre_type", "1"))
	diarySlice := service.NewBbsDiary().GetDiaryList(preType, dataId, userId, pageDepend)
	if len(diarySlice) < 16 {
		pageDepend++
		for _, v := range service.NewBbsDiary().GetDiaryList(preType, dataId, userId, pageDepend) {
			diarySlice = append(diarySlice, v)
		}
	}
	pageDepend++
	d.Json(http.StatusOK, "success", gin.H{"list": diarySlice, "pageDepend": pageDepend})
}

// 动态详情
func (d *diary) Detail() {
	userId := d.GetLoggedUserId()
	id, _ := strconv.Atoi(d.GetQueryByMethod("id", "0"))
	if id <= 0 {
		d.Json(http.StatusNotImplemented, "缺少必要参数id")
		return
	}
	if detail := service.NewBbsDiary().GetDiaryDetail(strconv.Itoa(id)); len(detail) > 0 {
		dtUserId, ok1 := detail["user_id"].(string)
		dtState, ok2 := detail["state"].(string)
		if (ok1 && dtUserId == strconv.Itoa(userId)) || (ok2 && dtState == "1") {
			d.Json(http.StatusOK, "success", detail)
			return
		}
	}
	d.Json(http.StatusOK, "success")
}

// 动态删除
func (d *diary) Delete() {
	userId := d.GetLoggedUserId()
	id, _ := strconv.Atoi(d.GetQueryByMethod("id", "0"))
	if id <= 0 {
		d.Json(http.StatusNotImplemented, "缺少必要参数id")
		return
	}

	var isMaster, isOwner int
	if diary := service.NewBbsDiary().GetDiaryDetail(strconv.Itoa(id)); len(diary) > 0 {
		if owner, ok := diary["user_id"].(string); ok && owner == strconv.Itoa(userId) {
			isOwner = 1
		} else {
			if dataId, ok := diary["data_id"].(string); ok {
				if circle := service.NewBbsCircle().CircleDetail(dataId); len(circle) > 0 {
					if cid, ok := circle["user_id"]; ok && cid == strconv.Itoa(userId) {
						isMaster = 1
					}
				}
			}
		}
	} else {
		d.Json(http.StatusNotImplemented, "动态不存在")
		return
	}

	if isMaster == 1 || isOwner == 1 {
		service.NewBbsDiary().DeleteDiary(strconv.Itoa(id))
		d.Json(http.StatusOK, "删除成功")
	} else {
		d.Json(http.StatusNotImplemented, "没有删除权限")
	}
}

// 是否可以删除
func (d *diary) Check() {
	userId := d.GetLoggedUserId()
	id, _ := strconv.Atoi(d.GetQueryByMethod("id", "0"))
	if id <= 0 {
		d.Json(http.StatusNotImplemented, "缺少必要参数id")
		return
	}

	var isMaster, isOwner int
	if diary := service.NewBbsDiary().GetDiaryDetail(strconv.Itoa(id)); len(diary) > 0 {
		if owner, ok := diary["user_id"].(string); ok && owner == strconv.Itoa(userId) {
			isOwner = 1
		} else {
			if dataId, ok := diary["data_id"].(string); ok {
				if circle := service.NewBbsCircle().CircleDetail(dataId); len(circle) > 0 {
					if cid, ok := circle["user_id"]; ok && cid == strconv.Itoa(userId) {
						isMaster = 1
					}
				}
			}
		}
	} else {
		d.Json(http.StatusNotImplemented, "动态不存在")
		return
	}

	if isMaster == 1 || isOwner == 1 {
		d.Json(http.StatusOK, "可以删除")
	} else {
		d.Json(http.StatusNotImplemented, "没有删除权限")
	}
}
