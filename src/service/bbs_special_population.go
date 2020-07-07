package service

import (
	"src/abs"
	app "src/config"
	"src/model"
	"strconv"
	"time"
)

type bbsSpecialPopulation struct {
	abs.Service
}

func NewBbsSpecialPopulation() *bbsSpecialPopulation {
	return &bbsSpecialPopulation{}
}

// 达人判断
func (sp *bbsSpecialPopulation) GetSpecialPopulationDetail(id string) map[string]string {
	var population map[string]string
	if jsonStr := sp.GetCache().HGet(app.Redis.Key.Bbs.HashSpecialPopulation, id).Val(); jsonStr == "" {
		resp := sp.ExamineSpecialPopulation([]string{id})
		population, _ = resp[id].(map[string]string)
	} else {
		_ = app.Json.Unmarshal([]byte(jsonStr), &population)
	}
	return population
}

// 审核达人
func (sp *bbsSpecialPopulation) ExamineSpecialPopulation(ids []string) map[string]interface{} {
	result := map[string]interface{}{}
	setMap := map[string]interface{}{}
	list := model.NewBbsSpecialPopulation().FindSpecialPopulation(ids, "id,user_id,state,circle_data")
	if len(list) > 0 {
		circleData := ""
		var state = "1"
		userIdStr := ""
		for _, v := range list {
			userIdStr = strconv.Itoa(v.UserId)
			if v.State == 1 {
				circleData = v.CircleData
				state = "1"
			}
			result[userIdStr] = map[string]string{"circle_data": circleData, "state": state}
			setMap[userIdStr], _ = app.Json.Marshal(result[userIdStr])
		}
	}

	if len(ids) > len(setMap) {
		for _, v := range ids {
			if _, ok := setMap[v]; ok == false {
				result[v] = map[string]string{"circle_data": "", "state": "0"}
				setMap[v], _ = app.Json.Marshal(result[v])
			}
		}
	}

	if len(setMap) > 0 {
		cache2 := sp.GetCache()
		_ = cache2.HMSet(app.Redis.Key.Bbs.HashSpecialPopulation, setMap).Err()
		if cache2.TTL(app.Redis.Key.Bbs.HashSpecialPopulation).Val() == -1*time.Second {
			cache2.Expire(app.Redis.Key.Bbs.HashSpecialPopulation, app.Redis.NormalTTL)
		}
	}
	return result
}
