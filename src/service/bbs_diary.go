package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	"net/http"
	"path"
	"src/abs"
	app "src/config"
	"src/library/tool"
	model "src/model"
	table "src/table"
	"strconv"
	"strings"
	"time"
)

type bbsDiary struct {
	abs.Service
}

func NewBbsDiary() *bbsDiary {
	return &bbsDiary{}
}

// 新增动态
func (dd *bbsDiary) AddDiary(diary table.BbsDiary) int {
	insertId := 0
	if insertId = model.NewBbsDiary().InsertDiary(diary); insertId > 0 {
		dd.ExamineDiary([]string{strconv.Itoa(insertId)})
	}
	return insertId
}

// 删除动态
func (dd *bbsDiary) DeleteDiary(id string) {
	intId, _ := strconv.Atoi(id)
	if num := model.NewBbsDiary().DeleteDiary(intId); num > 0 {
		dd.ExamineDiary([]string{id})
	}
}

// 动态详情
func (dd *bbsDiary) GetDiaryDetail(id string) map[string]interface{} {
	result := map[string]interface{}{}
	one := map[string]string{}
	jsonStr := dd.GetCache().HGet(app.Redis.Key.Bbs.HashDiary, tool.FixPad(id)).Val()
	if jsonStr == "" {
		resp := dd.ExamineDiary([]string{id})
		if v, ok := resp[id]; ok {
			one, _ = v.(map[string]string)
		}
	} else {
		_ = app.Json.Unmarshal([]byte(jsonStr), &one)
	}
	if len(one) > 0 {
		subData := map[string]interface{}{}
		_ = app.Json.Unmarshal([]byte(one["sub_data"]), &subData)
		result = map[string]interface{}{
			"content":     one["content"],
			"create_time": one["create_time"],
			"data_id":     one["data_id"],
			"data_type":   one["data_type"],
			"id":          one["id"],
			"state":       one["state"],
			"sub_data":    subData,
			"sub_type":    one["sub_type"],
			"user_id":     one["user_id"],
			"circle_name": one["circle_name"],
			"circle_icon": one["circle_icon"],
			"nickname":    "caius",
			"head_img":    "https://hxsapp-oss.hxsapp.com/public/image/default_head_img_new.png",
		}
		for _, v := range dd.GetUserInfo([]string{one["user_id"]}) {
			if one, ok := v.(map[string]interface{}); ok {
				result["nickname"] = one["nickname"]
				result["head_img"] = one["head_img"]
			}
		}
	}
	return result
}

// 动态列表
func (dd *bbsDiary) GetDiaryList(preType, dataId, userId, page int) []interface{} {
	var result []interface{}
	diaryMap := map[string]interface{}{}
	diarySlice := []string{}
	var ids []string
	ssd := dd.GetSSD()
	var pCount int64 = 20
	var pOffset int64 = 0
	pOffset = ((int64(page)) - 1) * pCount

	switch preType {
	case 1:
		// 圈子动态
		circleKey := app.Redis.Key.Bbs.ZSetDiaryCircle
		ownerKey := app.Redis.Key.Bbs.ZSetDiaryCircleOwner + strconv.Itoa(userId)
		unionKey := app.Redis.Key.Bbs.ZSetUnionDiary + strconv.Itoa(userId)
		if ssd.Exists(ownerKey).Val() > 0 {
			uKey := unionKey
			if _, err := ssd.ZUnionStore(unionKey, redis.ZStore{}, circleKey, ownerKey).Result(); err != nil {
				uKey = circleKey
			}
			diarySlice = ssd.ZRevRangeByScore(uKey, redis.ZRangeBy{Min: strconv.Itoa(dataId), Max: strconv.Itoa(dataId), Offset: pOffset, Count: pCount}).Val()
		} else {
			diarySlice = ssd.ZRevRangeByScore(circleKey, redis.ZRangeBy{Min: strconv.Itoa(dataId), Max: strconv.Itoa(dataId), Offset: pOffset, Count: pCount}).Val()
		}
	case 2:
		//　首页动态
		diaryKey := app.Redis.Key.Bbs.ZSetDiaryList
		diarySlice = ssd.ZRevRangeByScore(diaryKey, redis.ZRangeBy{Min: "-inf", Max: "+inf", Offset: pOffset, Count: pCount}).Val()
	default:
		// 个人主页动态
		diarySlice = ssd.ZRevRangeByScore(app.Redis.Key.Bbs.ZSetDiaryOwner,
			redis.ZRangeBy{Min: strconv.Itoa(dataId), Max: strconv.Itoa(dataId), Offset: pOffset, Count: pCount}).Val()
	}

	if len(diarySlice) > 0 {
		for k, v := range dd.GetCache().HMGet(app.Redis.Key.Bbs.HashDiary, diarySlice...).Val() {
			id := tool.FixTrim(diarySlice[k])
			one := map[string]string{}
			if v == nil {
				ids = append(ids, id)
				diaryMap[id] = map[string]string{}
			} else {
				if err := app.Json.Unmarshal([]byte(v.(string)), &one); err == nil {
					diaryMap[id] = one
				}
			}
		}
	}

	if len(ids) > 0 {
		for k, v := range dd.ExamineDiary(ids) {
			vMap, _ := v.(map[string]string)
			diaryMap[k] = vMap
		}
	}

	myCircle := NewBbsCircle().GetCircleIcon(userId, 1, 100)
	for _, v := range diarySlice {
		one, _ := diaryMap[tool.FixTrim(v)].(map[string]string)
		isMyVisibleCircle := 0
		if len(one) == 0 {
			continue
		}

		for _, c := range myCircle {
			if one["data_id"] == c {
				isMyVisibleCircle = 1
				break
			}
		}

		// 可见圈子
		if isMyVisibleCircle <= 0 {
			continue
		}

		// 待审的动态　本人
		if one["state"] == "0" && one["user_id"] != strconv.Itoa(userId) {
			continue
		}

		subData := map[string]interface{}{}
		_ = app.Json.Unmarshal([]byte(one["sub_data"]), &subData)
		result = append(result, map[string]interface{}{
			"id":           one["id"],
			"user_id":      one["user_id"],
			"content":      one["content"],
			"sub_type":     one["sub_type"],
			"sub_data":     subData,
			"data_id":      one["data_id"],
			"data_type":    one["data_type"],
			"state":        one["state"],
			"create_time":  one["create_time"],
			"circle_name":  one["circle_name"],
			"circle_icon":  one["circle_icon"],
			"nickname":     "caius",
			"head_img":     app.Bbs.HeadImageGirl,
			"is_commend":   "1",
			"commend_nums": "1",
		})
	}
	if len(result) > 0 {
		dd.dealData(result, userId)
	}
	return result
}

// 动态兼容用户行为
func (dd *bbsDiary) dealData(data []interface{}, userId int) {
	userInfo := map[string]interface{}{}
	diaryLike := map[string]interface{}{}
	circleInfo := map[string]interface{}{}
	var userIds []string
	var diaryIds []int
	var circleIds []string
	for _, v := range data {
		if one, ok := v.(map[string]interface{}); ok {
			if userId, ok := one["user_id"]; ok {
				userIds = append(userIds, userId.(string))
			}
			if id, ok := one["id"]; ok {
				idInt, _ := strconv.Atoi(id.(string))
				diaryIds = append(diaryIds, idInt)
			}
			if dataId, ok := one["data_id"]; ok {
				circleIds = append(circleIds, tool.FixPad(dataId.(string)))
			}
		}
	}
	if len(userIds) > 0 {
		for k, v := range dd.GetUserInfo(userIds) {
			if one, ok := v.(map[string]interface{}); ok {
				userInfo[k] = one
			}
		}
	}
	if len(circleIds) > 0 {
		for k, v := range NewBbsCircle().GetCircleInfo(circleIds) {
			if one, ok := v.(map[string]string); ok {
				circleInfo[k] = one
			}
		}
	}
	if len(diaryIds) > 0 {
		for k, v := range NewBbsDiaryLike().GetAllLikeCounter(diaryIds) {
			diaryLike[strconv.Itoa(k)] = map[string]string{"commend_nums": strconv.Itoa(v), "is_commend": "0"}
		}
		for k, v := range NewBbsDiaryLike().UserIsLike(userId, diaryIds) {
			if _, ok := diaryLike[strconv.Itoa(k)]; ok {
				if v {
					diaryLike[strconv.Itoa(k)].(map[string]string)["is_commend"] = "1"
				} else {
					diaryLike[strconv.Itoa(k)].(map[string]string)["is_commend"] = "0"
				}
			} else {
				diaryLike[strconv.Itoa(k)] = map[string]string{"commend_nums": "0", "is_commend": "1"}
			}
		}
	}

	for _, v := range data {
		if one, ok := v.(map[string]interface{}); ok {
			if u, ok := userInfo[one["user_id"].(string)]; ok {
				if uOne, ok := u.(map[string]interface{}); ok {
					one["nickname"] = uOne["nickname"]
					one["head_img"] = uOne["head_img"]
				}
			}
			if u, ok := diaryLike[one["id"].(string)]; ok {
				if uOne, ok := u.(map[string]string); ok {
					one["commend_nums"] = uOne["commend_nums"]
					one["is_commend"] = uOne["is_commend"]
				}
			}
			if u, ok := circleInfo[one["data_id"].(string)]; ok {
				if uOne, ok := u.(map[string]string); ok {
					one["circle_name"] = uOne["title"]
					one["circle_icon"] = uOne["icon"]
				}
			}
		}
	}

}

// 审核动态
func (dd *bbsDiary) ExamineDiary(ids []string) map[string]interface{} {
	zRemSliceMap := map[string][]interface{}{}
	var ZSetUnionSlice []string
	zSetAddMap := map[string][]redis.Z{}
	setMap := map[string]interface{}{}
	result := map[string]interface{}{}
	for _, v := range model.NewBbsDiary().FindDiary(ids, "*") {
		fId := tool.FixPad(strconv.Itoa(v.Id))
		userIdFloat64, _ := strconv.ParseFloat(strconv.Itoa(v.UserId), 64)
		dataIdFloat64, _ := strconv.ParseFloat(strconv.Itoa(v.DataId), 64)
		idFloat64, _ := strconv.ParseFloat(strconv.Itoa(v.Id), 64)
		userIdStr := strconv.Itoa(v.UserId)
		one := map[string]string{
			"id":          strconv.Itoa(v.Id),
			"user_id":     strconv.Itoa(v.UserId),
			"content":     v.Content,
			"sub_data":    v.SubData,
			"sub_type":    strconv.Itoa(v.SubType),
			"data_id":     strconv.Itoa(v.DataId),
			"data_type":   strconv.Itoa(v.DataType),
			"state":       strconv.Itoa(v.State),
			"create_time": tool.ParseTime(v.CreateTime, "2006-01-02T15:04:05+08:00"),
			"circle_name": v.CircleName,
			"circle_icon": v.Icon,
		}
		subData := map[string]interface{}{}
		var dataSlice []string
		_ = app.Json.Unmarshal([]byte(one["sub_data"]), &dataSlice)
		if len(dataSlice) > 0 {
			for k, url := range dataSlice {
				dataSlice[k] = tool.DealPreOss(url, v.SubType, 1)
			}
			if one["sub_type"] == "1" {
				imageInfo := map[string]string{"width": "0", "height": "0"}
				imageInfo["height"], imageInfo["width"] = dd.GetImageInfo(dataSlice[0])
				subData = map[string]interface{}{"image": dataSlice, "size": imageInfo}
			} else {
				snapshot := ""
				if len(dataSlice) > 0 {
					fS := path.Base(dataSlice[0])         //获取文件名带后缀
					fE := path.Ext(dataSlice[0])          //获取文件后缀
					snapshot = strings.TrimSuffix(fS, fE) //获取文件名
					snapshot = fmt.Sprintf(app.Bbs.VideoSnapshot, snapshot)
				}
				subData = map[string]interface{}{"video": dataSlice, "snapshot": snapshot}
			}
		}
		if len(subData) > 0 {
			jsonSubData, _ := app.Json.Marshal(subData)
			one["sub_data"] = string(jsonSubData)
		}
		setMap[fId], _ = app.Json.Marshal(one)
		result[one["id"]] = one
		switch v.State {
		case 0:
			if _, ok := zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircle]; ok {
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircle] = append(zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircle], fId)
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryList] = append(zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryList], fId)
			} else {
				// 圈子审核通过的动态
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircle] = []interface{}{fId}
				// 社区首页的动态
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryList] = []interface{}{fId}
			}
			if _, ok := zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr]; ok {
				zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr] = append(
					zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr],
					redis.Z{Score: dataIdFloat64, Member: fId})
			} else {
				zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr] = []redis.Z{{dataIdFloat64, fId}}
			}
			if _, ok := zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryOwner]; ok {
				zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryOwner] = append(zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryOwner], redis.Z{Score: dataIdFloat64, Member: fId})

			} else {
				zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryOwner] = []redis.Z{{userIdFloat64, fId}}
			}
		case 1:
			if _, ok := zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryCircle]; ok {
				zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryCircle] = append(zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryCircle], redis.Z{Score: dataIdFloat64, Member: fId})
				zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryOwner] = append(zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryOwner], redis.Z{Score: userIdFloat64, Member: fId})
				zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryList] = append(zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryList], redis.Z{Score: idFloat64, Member: fId})
			} else {
				zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryCircle] = []redis.Z{{dataIdFloat64, fId}}
				zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryOwner] = []redis.Z{{userIdFloat64, fId}}
				zSetAddMap[app.Redis.Key.Bbs.ZSetDiaryList] = []redis.Z{{idFloat64, fId}}
			}
			if _, ok := zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr]; ok {
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr] = append(zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr], fId)
			} else {
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr] = []interface{}{fId}
			}
		default:
			if _, ok := zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr]; ok {
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr] = append(zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr], fId)
			} else {
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircleOwner+userIdStr] = []interface{}{fId}
			}
			if _, ok := zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircle]; ok {
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircle] = append(zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircle], fId)
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryOwner] = append(zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryOwner], fId)
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryList] = append(zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryList], fId)
			} else {
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryCircle] = []interface{}{fId}
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryOwner] = []interface{}{fId}
				zRemSliceMap[app.Redis.Key.Bbs.ZSetDiaryList] = []interface{}{fId}
			}
		}
		ZSetUnionSlice = append(ZSetUnionSlice, app.Redis.Key.Bbs.ZSetUnionDiary+userIdStr)
	}
	ssd := dd.GetSSD()
	for k, v := range zRemSliceMap {
		_ = ssd.ZRem(k, v...)
	}
	for k, v := range zSetAddMap {
		_ = ssd.ZAdd(k, v...)
	}
	if len(ZSetUnionSlice) > 0 {
		_ = ssd.Del(ZSetUnionSlice...)
	}
	if len(setMap) > 0 {
		cache := dd.GetCache()
		_ = cache.HMSet(app.Redis.Key.Bbs.HashDiary, setMap).Err()
		if cache.TTL(app.Redis.Key.Bbs.HashDiary).Val() == -1*time.Second {
			cache.Expire(app.Redis.Key.Bbs.HashDiary, app.Redis.NormalTTL)
		}
	}
	return result
}

// 图片的尺寸
func (dd *bbsDiary) GetImageInfo(url string) (string, string) {
	imgH, imgW := "0", "0"
	res, err := tool.NewRequest(url + "?x-oss-process=image/info").SetTimeOut(2 * time.Second).Get()
	if err == nil && res.Body != "" {
		if gjson.Get(res.Body, "ImageHeight.value").Exists() {
			imgH = gjson.Get(res.Body, "ImageHeight.value").String()
		}
		if gjson.Get(res.Body, "ImageWidth.value").Exists() {
			imgW = gjson.Get(res.Body, "ImageWidth.value").String()
		}
	}
	return imgH, imgW
}

// 用户信息
func (dd *bbsDiary) GetUserInfo(ids []string) map[string]interface{} {
	result := map[string]interface{}{}
	params := map[string]interface{}{"ids": strings.Join(ids, ",")}
	res := tool.HttpRequestWithSign(app.Sdk.Account.GetShortUserInfo, params, http.MethodPost, 3*time.Second)
	if gjson.Get(res, "code").String() == "200" {
		if gjson.Get(res, "data.list").IsObject() {
			if v, ok := gjson.Get(res, "data.list").Value().(map[string]interface{}); ok {
				result = v
			}
		}
	}
	return result
}
