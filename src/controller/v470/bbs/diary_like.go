package v470

import (
	"net/http"
	"src/abs"
	"src/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// diaryLike 点赞模块
type diaryLike struct {
	abs.Controller
}

// DiaryComment 构造函数
func DiaryLike(c *gin.Context) *diaryLike {
	l := &diaryLike{}
	l.Ctx = c
	return l
}

// Save 点赞 / 取消点赞
func (d *diaryLike) Save() {
	userId := d.GetLoggedUserId()
	diaryId, _ := strconv.Atoi(d.GetQueryByMethod("diary_id", "0"))
	if diaryId <= 0 {
		d.Json(http.StatusNotImplemented, "缺少必要参数 diary_id")
		return
	}
	isLike := false
	if d.GetQueryByMethod("is_like", "0") == "1" {
		isLike = true
	}
	if result := service.NewBbsDiaryLike().Save(userId, isLike, diaryId); result {
		d.Json(http.StatusOK, "操作成功")
	} else {
		d.Json(511, "操作失败")
	}
}

// List 指定动态的点赞列表
func (d *diaryLike) List() {
	userId := d.GetLoggedUserId()
	diaryId, _ := strconv.Atoi(d.GetQueryByMethod("diary_id", "0"))
	if diaryId <= 0 {
		d.Json(http.StatusNotImplemented, "缺少必要参数 diary_id")
		return
	}
	dl := service.NewBbsDiaryLike()
	result := dl.GetLastLike(diaryId)
	counter := dl.GetAllLikeCounter([]int{diaryId})
	isLike := dl.UserIsLike(userId, []int{diaryId})
	d.Json(http.StatusOK, "success", gin.H{"list": result, "counter": counter[diaryId], "is_like": isLike[diaryId]})
}
