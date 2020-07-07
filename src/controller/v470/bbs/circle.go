package v470

import (
	"net/http"
	"src/abs"
	"src/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type circle struct {
	abs.Controller
}

func Circle(c *gin.Context) *circle {
	dc := &circle{}
	dc.Ctx = c
	return dc
}

// 圈子动态列表
func (c *circle) Icon() {
	userId := c.GetLoggedUserId()
	list := service.NewBbsCircle().GetCircleIcon(userId, 0, 8)
	c.Json(http.StatusOK, "success", gin.H{"list": list})
}

// 圈子详情
func (c *circle) Detail() {
	userId := c.GetLoggedUserId()
	id, _ := strconv.Atoi(c.GetQueryByMethod("id", "0"))
	if id <= 0 {
		c.Json(http.StatusNotImplemented, "缺少必要参数圈子id")
		return
	}
	one := service.NewBbsCircle().CircleDetail(strconv.Itoa(id))
	if _, ok := one["visible_user"]; ok {
		if service.NewBbsCircle().IsShow(strconv.Itoa(userId), one["user_id"], one["visible_user"], one["id"]) {
			delete(one, "sort")
			delete(one, "state")
			delete(one, "visible_user")
		} else {
			c.Json(http.StatusOK, "圈子不可见")
			return
		}
	} else {
		c.Json(http.StatusOK, "圈子不存在")
		return
	}
	c.Json(http.StatusOK, "success", one)

}

// 审核动态 更新缓存
func (c *circle) Examine() {
	if ids := c.GetQueryByMethod("ids", ""); ids != "" {
		_ = service.NewBbsCircle().ExamineCircle(strings.Split(ids, ","))
	}
	c.Json(http.StatusOK, "更新缓存成功")
}
