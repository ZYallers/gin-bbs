package v470

import (
	"net/http"
	"src/abs"
	"src/model"
	"src/service"
	"src/table"
	"strconv"

	"github.com/gin-gonic/gin"
)

type diaryComment struct {
	abs.Controller
}

// DiaryComment 构造函数
func DiaryComment(c *gin.Context) *diaryComment {
	l := &diaryComment{}
	l.Ctx = c
	return l
}

// Add 新增评论
func (d *diaryComment) Add() {
	userId := d.GetLoggedUserId()
	diaryId, _ := strconv.Atoi(d.GetQueryByMethod("diary_id", "0"))
	templateId, _ := strconv.Atoi(d.GetQueryByMethod("template_id", "0"))
	templateId = templateId - 1 // 数组索引是从 0 开始, tempate_id 从 1 开始
	var content string
	if diaryId <= 0 {
		d.Json(http.StatusNotImplemented, "缺少必要参数diary_id")
		return
	}

	if template := model.NewBbsDiaryComment().Template(); templateId < 0 || templateId >= len(template) {
		d.Json(http.StatusNotImplemented, "非法template_id")
		return
	} else {
		content = template[templateId]
	}
	
	if detail := service.NewBbsDiary().GetDiaryDetail(strconv.Itoa(diaryId)); len(detail) == 0 || detail[`state`] != `1` {
		d.Json(http.StatusNotImplemented, "动态正在审核中, 请稍候再来评论")
		return
	}

	comment := table.BbsDiaryComment{
		UserId:  userId,
		DiaryId: diaryId,
		Content: content,
		State:   model.NewBbsDiaryComment().DefaultState(),
	}
	if id := service.NewBbsDiaryComment().SaveComment(comment); id > 0 {
		d.Json(http.StatusOK, "发布成功")
	} else {
		d.Json(511, "发布失败")
	}
}

// List 评论列表
func (d *diaryComment) List() {
	diaryId, _ := strconv.Atoi(d.GetQueryByMethod("diary_id", "0"))
	if diaryId <= 0 {
		d.Json(http.StatusNotImplemented, "缺少必要参数 diary_id")
		return
	}
	page, _ := strconv.Atoi(d.GetQueryByMethod("page", "1"))
	size, _ := strconv.Atoi(d.GetQueryByMethod("size", "10"))
	dc := service.NewBbsDiaryComment()
	result := dc.GetAllComment(diaryId, page, size)
	counter := dc.GetCommentCounter(diaryId)
	d.Json(http.StatusOK, "success", gin.H{"list": result, "counter": counter})
}

// Template 评论模板
func (d *diaryComment) Template() {
	template := model.NewBbsDiaryComment().Template()
	result := make([]map[string]interface{}, len(template))
	for k, v := range template {
		result[k] = map[string]interface{}{`template_id`: k + 1, `content`: v}
	}
	d.Json(http.StatusOK, "success", gin.H{"list": result})
}
