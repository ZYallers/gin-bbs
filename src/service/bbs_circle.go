package service

import (
	"github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	"src/abs"
	app "src/config"
	"src/library/tool"
	"src/model"
	"strconv"
	"strings"
	"time"
)

type bbsCircle struct {
	abs.Service
	visibleUser int // 会员等级可见用户
}

func NewBbsCircle() *bbsCircle {
	return &bbsCircle{}
}

// 首页圈子icon
func (c *bbsCircle) GetCircleIcon(userId, short, num int) []interface{} {
	var (
		result                          []interface{}
		circleSlice, specialCircle, ids []string
	)
	circleMap := map[string]interface{}{}
	specialPopulation := NewBbsSpecialPopulation().GetSpecialPopulationDetail(strconv.Itoa(userId))
	if circleStr, ok := specialPopulation["circle_data"]; ok {
		for _, v := range strings.Split(circleStr, ",") {
			if vInt, _ := strconv.Atoi(v); vInt <= 0 {
				continue
			}
			specialCircle = append(specialCircle, tool.FixPad(v))
		}
	}

	circleSlice = c.GetSSD().ZRevRangeByScore(app.Redis.Key.Bbs.ZSetCircleList, redis.ZRangeBy{Min: "-inf", Max: "+inf", Count: 100}).Val()
	if len(circleSlice) > 0 {
		if len(specialCircle) > 0 {
			circleSlice = tool.UnionSlice(specialCircle, circleSlice)
		}
		for k, v := range c.GetCache().HMGet(app.Redis.Key.Bbs.HashCircle, circleSlice...).Val() {
			id := tool.FixTrim(circleSlice[k])
			if v == nil {
				ids = append(ids, id)
				circleMap[id] = map[string]string{}
			} else {
				one := map[string]string{}
				if err := app.Json.Unmarshal([]byte(v.(string)), &one); err == nil {
					if c.IsShow(strconv.Itoa(userId), one["user_id"], one["visible_user"], id) {
						circleMap[id] = one
					}
				}
			}
		}
	}
	if len(ids) > 0 {
		for k, v := range c.ExamineCircle(ids) {
			vMap, _ := v.(map[string]string)
			if c.IsShow(strconv.Itoa(userId), vMap["user_id"], vMap["visible_user"], vMap["id"]) {
				circleMap[k] = vMap
			}
		}
	}
	iN := 1
	for _, v := range circleSlice {
		key := tool.FixTrim(v)
		one, _ := circleMap[key].(map[string]string)
		if len(one) == 0 {
			continue
		}
		if short > 0 {
			result = append(result, one["id"])
		} else {
			if one["state"] != "1" {
				continue
			}
			result = append(result, map[string]string{"icon": one["icon"], "id": one["id"], "title": one["title"]})
		}
		if iN >= num {
			break
		}
		iN++
	}
	return result
}

// 判断圈子是否显示
func (c *bbsCircle) IsShow(userId, owner, visible, circleId string) bool {
	result := false
	visibleUser := c.GetUserVisible(userId)
	visibleCircle, _ := strconv.Atoi(visible)
	if userId == owner || (visibleUser&visibleCircle) > 0 {
		result = true
	} else {
		specialPopulation := NewBbsSpecialPopulation().GetSpecialPopulationDetail(userId)
		if circleStr, ok := specialPopulation["circle_data"]; ok {
			for _, v := range strings.Split(circleStr, ",") {
				if circleId == v {
					result = true
					break
				}
			}
		}
	}
	return result
}

// 审核圈子
func (c *bbsCircle) ExamineCircle(ids []string) map[string]interface{} {
	setMap := make(map[string]interface{}, len(ids))
	result := make(map[string]interface{}, len(ids))
	var zRemSlice []string
	zSetAddMap := map[string][]redis.Z{}
	for _, v := range model.NewBbsCircle().FindCircle(ids, "*") {
		fId := tool.FixPad(strconv.Itoa(v.Id))
		idStr := strconv.Itoa(v.Id)
		one := map[string]string{
			"id":           strconv.Itoa(v.Id),
			"user_id":      strconv.Itoa(v.UserId),
			"title":        v.Title,
			"introduction": v.Introduction,
			"icon":         v.Icon,
			"image":        v.Image,
			"sort":         strconv.Itoa(v.Sort),
			"state":        strconv.Itoa(v.State),
			"visible_user": strconv.Itoa(v.VisibleUser),
			"create_time":  tool.ParseTime(v.CreateTime, "2006-01-02T15:04:05+08:00"),
		}
		setMap[fId], _ = app.Json.Marshal(one)
		result[idStr] = one
		sortFloat64, _ := strconv.ParseFloat(strconv.Itoa(v.Sort), 64)
		if v.State > 0 {
			if _, ok := zSetAddMap[app.Redis.Key.Bbs.ZSetCircleList]; ok {
				zSetAddMap[app.Redis.Key.Bbs.ZSetCircleList] = append(zSetAddMap[app.Redis.Key.Bbs.ZSetCircleList], redis.Z{Score: sortFloat64, Member: fId})
			} else {
				zSetAddMap[app.Redis.Key.Bbs.ZSetCircleList] = []redis.Z{{Score: sortFloat64, Member: fId}}
			}
			for _, diaryId := range model.NewBbsDiary().FindCircleDiary(one["id"]) {
				idFloat64, _ := strconv.ParseFloat(diaryId, 64)
				if _, ok := zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryList]; ok {
					zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryList] = append(zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryList], redis.Z{Score: idFloat64, Member: diaryId})
				} else {
					zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryList] = []redis.Z{{Score: idFloat64, Member: diaryId}}
				}
			}
		} else {
			for _, diaryId := range model.NewBbsDiary().FindCircleDiary(one["id"]) {
				zRemSlice = append(zRemSlice, diaryId)
			}
		}
	}

	if len(setMap) > 0 {
		_ = c.GetCache().HMSet(app.Redis.Key.Bbs.HashCircle, setMap).Err()
		for k, v := range zSetAddMap {
			_ = c.GetSSD().ZAdd(k, v...).Err()
		}
		if c.GetCache().TTL(app.Redis.Key.Bbs.HashCircle).Val() == -1*time.Second {
			c.GetCache().Expire(app.Redis.Key.Bbs.HashCircle, app.Redis.NormalTTL)
		}
		if len(zRemSlice) > 0 {
			_ = c.GetSSD().ZRem(app.Redis.Key.Bbs.ZSetDiaryList, zRemSlice).Err()
		}
	}
	return result
}

// 圈子详情
func (c *bbsCircle) CircleDetail(id string) map[string]string {
	result := map[string]string{}
	jsonStr := c.GetCache().HGet(app.Redis.Key.Bbs.HashCircle, tool.FixPad(id)).Val()
	if jsonStr == "" {
		resp := c.ExamineCircle([]string{id})
		result, _ = resp[id].(map[string]string)
	} else {
		_ = app.Json.Unmarshal([]byte(jsonStr), &result)
	}
	stateInt, _ := strconv.Atoi(result["state"])
	if stateInt < 1 {
		result = map[string]string{}
	} else {
		for _, v := range NewBbsDiary().GetUserInfo([]string{result["user_id"]}) {
			if one, ok := v.(map[string]interface{}); ok {
				result["nickname"] = one["nickname"].(string)
				result["head_img"] = one["head_img"].(string)
			}
		}
	}
	return result
}

// 查询圈子
func (c *bbsCircle) GetCircleInfo(circleSlice []string) map[string]interface{} {
	ids := []string{}
	circleMap := map[string]interface{}{}
	for k, v := range c.GetCache().HMGet(app.Redis.Key.Bbs.HashCircle, circleSlice...).Val() {
		id := tool.FixTrim(circleSlice[k])
		if v == nil {
			ids = append(ids, id)
			circleMap[id] = map[string]string{}
		} else {
			one := map[string]string{}
			if err := app.Json.Unmarshal([]byte(v.(string)), &one); err == nil {
				circleMap[id] = one
			}
		}
	}

	if len(ids) > 0 {
		for k, v := range c.ExamineCircle(ids) {
			vMap, _ := v.(map[string]string)
			circleMap[k] = vMap
		}
	}
	return circleMap
}

// 用户等级
func (c *bbsCircle) GetUserVisible(userId string) int {
	if c.visibleUser <= 0 {
		params := map[string]interface{}{"user_id": userId}
		res := tool.HttpRequestWithSign(app.Sdk.Account.GetUserVipVisible, params, "POST", time.Second*3)
		if gjson.Get(res, "code").String() == "200" {
			if gjson.Get(res, "data.visible").Value() != nil {
				c.visibleUser, _ = strconv.Atoi(gjson.Get(res, "data.visible").String())
			}
		}
	}
	return c.visibleUser
}
