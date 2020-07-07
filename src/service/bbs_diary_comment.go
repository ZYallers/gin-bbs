package service

import (
	"fmt"
	"go.uber.org/zap"
	"src/abs"
	app "src/config"
	"src/model"
	"src/table"
	"strconv"
	"time"
)

type bbsDiaryComment struct {
	abs.Service
}

func NewBbsDiaryComment() *bbsDiaryComment {
	return &bbsDiaryComment{}
}

// SaveComment 添加评论, 更新计数器, 清除相关缓存
func (dc *bbsDiaryComment) SaveComment(comment table.BbsDiaryComment) int {
	insertId := 0
	if id := model.NewBbsDiaryComment().InsertComment(comment); id > 0 {
		dc.DeleteCache(comment.DiaryId) // 删除 diaryId 下的评论缓存
		insertId = id
	}
	return insertId
}

// DeleteCache 删除 diaryId 下的评论缓存
func (dc *bbsDiaryComment) DeleteCache(diaryId int) map[string]int64 {
	cache := dc.GetCache()
	listKey := fmt.Sprintf(app.Redis.Key.Bbs.DiaryCommentHash, diaryId)
	counterKey := fmt.Sprintf(app.Redis.Key.Bbs.DiaryCommentCounterString, diaryId)
	return map[string]int64{`list`: cache.Del(listKey).Val(), `counter`: cache.Del(counterKey).Val()}
}

// GetAllComment 获取上线的评论列表
func (dc *bbsDiaryComment) GetAllComment(diaryId int, page int, size int) (result []map[string]interface{}) {
	cache := dc.GetCache()
	listKey := fmt.Sprintf(app.Redis.Key.Bbs.DiaryCommentHash, diaryId)
	field := fmt.Sprintf(`%d_%d`, page, size)
	cacheData := cache.HGet(listKey, field).Val()
	if cacheData == `` {
		commentList := model.NewBbsDiaryComment().FindComment(diaryId, (page-1)*size, size) // 获取上线的评论列表
		if len(commentList) > 0 {
			for _, v := range commentList {
				result = append(result, map[string]interface{}{
					`id`:          v.Id,
					`diary_id`:    v.DiaryId,
					`user_id`:     v.UserId,
					`content`:     v.Content,
					`create_time`: v.CreateTime.Format("2006/01/02 15:04"),
				})
			}
			encode, err := app.Json.Marshal(&result)
			if err == nil {
				cache.HSet(listKey, field, encode)
			} else {
				app.Logger.Error(`GetAllComment 方法中 json 编码失败`, zap.Any(`result`, result))
			}
		} else {
			cache.HSet(listKey, field, []byte(`-1`)) // -1=查询数据库的结果为空
		}

		ttl := cache.TTL(listKey).Val()
		if ttl == -1*time.Second {
			cache.Expire(listKey, 60*5*time.Second)
		}
	} else if cacheData != `-1` { // -1=查询数据库的结果为空, 直接返回空的 result
		cacheData = cache.HGet(listKey, field).Val()
		err := app.Json.Unmarshal([]byte(cacheData), &result)
		if err != nil {
			app.Logger.Error(`GetAllComment 方法中 json 解码失败`, zap.Any(`cacheData`, cacheData))
		} else {
			for k, v := range result { // 把 user_id 从 float64 转换为 int
				result[k][`user_id`] = int(v[`user_id`].(float64))
			}
		}
	}
	if len(result) > 0 { // 批量追加用户昵称和头像地址
		result = NewUserInfo().AppendInfo(result, []string{`head_img`, `nickname`})
	}
	return result
}

// GetCommentCounter 获取评论计数器
func (dc *bbsDiaryComment) GetCommentCounter(diaryId int) (result int) {
	cache := dc.GetCache()
	counterKey := fmt.Sprintf(app.Redis.Key.Bbs.DiaryCommentCounterString, diaryId)
	cacheData := cache.Get(counterKey).Val()
	if cacheData == `` {
		result = model.NewBbsDiaryComment().GetCounter(diaryId) // 获取评论计数器
		cache.Set(counterKey, result, 60*60*time.Second)
	} else {
		result, _ = strconv.Atoi(cacheData)
	}
	return result
}
