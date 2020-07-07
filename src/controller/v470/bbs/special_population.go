package v470

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"src/abs"
	"src/service"
	"strconv"
	"strings"
)

type specialPopulation struct {
	abs.Controller
}

func SpecialPopulation(c *gin.Context) *specialPopulation {
	sp := &specialPopulation{}
	sp.Ctx = c
	return sp
}

// 达人详情
func (sp *specialPopulation) Detail() {
	info := map[string]interface{}{"is_special_population": 0}
	userId, _ := strconv.Atoi(sp.GetQueryByMethod("user_id", "0"))
	resp := service.NewBbsSpecialPopulation().GetSpecialPopulationDetail(strconv.Itoa(userId))
	if v, ok := resp["state"]; ok && v == "1" {
		info["is_special_population"] = 1
	}
	sp.Json(http.StatusOK, "success", info)
}

// 审核动态 更新缓存
func (sp *specialPopulation) Examine() {
	if ids := sp.GetQueryByMethod("ids", ""); ids != "" {
		_ = service.NewBbsSpecialPopulation().ExamineSpecialPopulation(strings.Split(ids, ","))
	}
	sp.Json(http.StatusOK, "更新缓存成功")
}
