package service

import (
	"fmt"
	"src/abs"
	app "src/config"
	"src/model"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type bbsDiaryLike struct {
	abs.Service
}

func NewBbsDiaryLike() *bbsDiaryLike {
	return &bbsDiaryLike{}
}

// Save 点赞 / 取消点赞, isLike=true 代表点赞， isLike=false 代表取消点赞
func (dl *bbsDiaryLike) Save(userId int, isLike bool, diaryId int) (result bool) {
	if result = (&model.BbsDiaryLike{}).Save(userId, isLike, diaryId); result {
		dl.DeleteCache(userId, diaryId) // 删除 diaryId 下的点赞缓存
	}
	return result
}

// DeleteCache 删除 diaryId 下的点赞缓存
func (dl *bbsDiaryLike) DeleteCache(userId int, diaryId int) map[string]int64 {
	cache := dl.GetCache()
	listKey := fmt.Sprintf(app.Redis.Key.Bbs.DiaryLikeZSet, diaryId)
	counterKey := fmt.Sprintf(app.Redis.Key.Bbs.DiaryLikeCounterString, diaryId)
	userLikeKey := fmt.Sprintf(app.Redis.Key.Bbs.UserLikeDiaryString, userId, diaryId)
	return map[string]int64{
		`list`:    cache.Del(listKey).Val(),
		`counter`: cache.Del(counterKey).Val(),
		`user`:    cache.Del(userLikeKey).Val(),
	}
}

// GetLastLike 获取最新7个点赞数据
func (dl *bbsDiaryLike) GetLastLike(diaryId int) (result []map[string]interface{}) {
	var num int64 = 7 // 需求是最新7个点赞数据
	cache := dl.GetCache()
	listKey := fmt.Sprintf(app.Redis.Key.Bbs.DiaryLikeZSet, diaryId)
	cacheData := cache.ZRevRangeByScoreWithScores(listKey, redis.ZRangeBy{Min: `-inf`, Max: `+inf`, Offset: 0, Count: num}).Val()
	if len(cacheData) == 0 {
		likeList := (&model.BbsDiaryLike{}).FindLastLike(diaryId, num) // 获取最新 num 个点赞数据
		if len(likeList) > 0 {
			for _, v := range likeList {
				unixTime, _ := strconv.ParseFloat(strconv.FormatInt(v.CreateTime.Unix(), 10), 64)
				userId := strconv.Itoa(v.UserId)
				cache.ZAdd(listKey, redis.Z{Score: unixTime, Member: userId})
				result = append(result, map[string]interface{}{`create_time`: unixTime, `user_id`: userId})
			}
		} else {
			cache.ZAdd(listKey, redis.Z{Score: -1, Member: `-1`}) // -1=查询数据库的结果为空
		}

		ttl := cache.TTL(listKey).Val()
		if ttl == -1*time.Second {
			cache.Expire(listKey, 60*5*time.Second)
		}
	} else if cacheData[0].Score != -1 { // -1=查询数据库的结果为空
		for _, v := range cacheData {
			result = append(result, map[string]interface{}{
				`create_time`: v.Score,
				`user_id`:     v.Member,
			})
		}
	}
	if len(result) > 0 { // 批量追加用户头像地址
		for k, v := range result { // 把 user_id 从 string 转换为 int
			result[k][`user_id`], _ = strconv.Atoi(v[`user_id`].(string))
		}
		result = NewUserInfo().AppendInfo(result, []string{`head_img`})
	}
	return result
}

// GetLikeCounter 批量获取点赞计数器
func (dl *bbsDiaryLike) GetAllLikeCounter(diaryId []int) (result map[int]int) {
	result = make(map[int]int)
	if len(diaryId) > 0 {
		cache := dl.GetCache()
		var counterKey []string
		var noCacheDiaryId []int
		for _, v := range diaryId { // 使用 slice 包起这个用户所有的 key
			counterKey = append(counterKey, fmt.Sprintf(app.Redis.Key.Bbs.DiaryLikeCounterString, v))
		}
		cacheData := cache.MGet(counterKey...).Val()
		for k, v := range cacheData {
			if v != nil { // 缓存里面有值
				result[diaryId[k]], _ = strconv.Atoi(v.(string))
			} else { // 缓存没有值, 放入查询数据库的变量中
				noCacheDiaryId = append(noCacheDiaryId, diaryId[k])
			}
		}
		if len(noCacheDiaryId) > 0 {
			var noCacheData []interface{}
			dbResult := (&model.BbsDiaryLike{}).FindAllCounter(noCacheDiaryId) // 批量从数据库中检查用户是否对指定的动态点赞
			for k, v := range dbResult {
				noCacheData = append(noCacheData, k)
				noCacheData = append(noCacheData, v)
				result[k] = v
			}
			cache.MSet(noCacheData...).Val()
		}
	}
	return result
}

// UserIsLike 批量判断用户是否对指定 diaryId 点赞过
func (dl *bbsDiaryLike) UserIsLike(userId int, diaryId []int) (result map[int]bool) {
	result = make(map[int]bool)
	if len(diaryId) > 0 {
		cache := dl.GetCache()
		var userLikeKey []string
		var noCacheDiaryId []int
		for _, v := range diaryId { // 使用 slice 包起这个用户所有的 key
			userLikeKey = append(userLikeKey, fmt.Sprintf(app.Redis.Key.Bbs.UserLikeDiaryString, userId, v))
		}
		cacheData := cache.MGet(userLikeKey...).Val()
		for k, v := range cacheData {
			if v != nil { // 缓存里面有值
				result[diaryId[k]] = v.(string) == `-1`
			} else { // 缓存没有值, 放入查询数据库的变量中
				noCacheDiaryId = append(noCacheDiaryId, diaryId[k])
			}
		}
		if len(noCacheDiaryId) > 0 {
			var noCacheData []interface{}
			dbResult := (&model.BbsDiaryLike{}).FindUserLike(userId, noCacheDiaryId) // 批量从数据库中检查用户是否对指定的动态点赞
			for k, v := range dbResult {
				noCacheData = append(noCacheData, k)
				if v == true {
					noCacheData = append(noCacheData, `1`)
				} else {
					noCacheData = append(noCacheData, `-1`)
				}
				result[k] = v
			}
			cache.MSet(noCacheData...).Val()
		}
	}
	return result
}
