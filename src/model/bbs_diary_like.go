package model

import (
	"go.uber.org/zap"
	"src/config"
	"src/table"
)

type BbsDiaryLike struct {
	baseModel
}

// Save 点赞 / 取消点赞, isLike=true 代表点赞， isLike=false 代表取消点赞
func (dl *BbsDiaryLike) Save(userId int, isLike bool, diaryId int) (result bool) {
	row := table.BbsDiaryLike{DiaryId: diaryId, UserId: userId}
	var err error
	if isLike {
		row.State = dl.DefaultState()
		err = dl.GetBbs().Create(&row).Error
	} else {
		err = dl.GetBbs().Unscoped().Delete(&row, `user_id = ? AND diary_id = ?`, userId, diaryId).Error
	}
	if err != nil {
		app.Logger.Error(`BbsDiaryLike.Save() 失败 `,
			zap.Int(`diaryId`, diaryId),
			zap.Int(`userId`, userId),
			zap.Bool(`isLike`, isLike),
			zap.NamedError(`err`, err))
	} else {
		result = true
	}
	return
}

// FindLastLike 返回最新 num 个点赞数据
func (dl *BbsDiaryLike) FindLastLike(diaryId int, num int64) (result []table.BbsDiaryLike) {
	dl.GetBbs().Select(`id,diary_id,user_id,create_time`).
		Where(`diary_id = ? AND state = ?`, diaryId, StateOnline).
		Order(`id DESC`).Limit(num).Find(&result)
	return result
}

// GetCounter 批量获取上线的点赞总数
func (dl *BbsDiaryLike) FindAllCounter(diaryId []int) (result map[int]int) {
	result = make(map[int]int, len(diaryId))
	if len(diaryId) > 0 {
		var total, id int
		rows, _ := dl.GetBbs().Model(&table.BbsDiaryLike{}).Select(`diary_id, COUNT(id) AS total`).
			Where(`diary_id IN (?) AND state = ?`, diaryId, StateOnline).Group(`diary_id`).Rows()
		for rows.Next() {
			_ = rows.Scan(&id, &total)
			result[id] = total
		}

		for _, v := range diaryId {
			if _, ok := result[v]; !ok { // 数据不存在的时候默认返回 0
				result[v] = 0
			}
		}
	}
	return
}

// FindUserLike 批量检查用户是否对指定的动态点赞
// userId 用户 id
// diaryId 动态 id
// result 返回结果为 map[diaryId]bool
func (dl *BbsDiaryLike) FindUserLike(userId int, diaryId []int) map[int]bool {
	result := make(map[int]bool, len(diaryId))
	var data []table.BbsDiaryLike
	dl.GetBbs().Model(&table.BbsDiaryLike{}).Select(`diary_id`).
		Where(`user_id =? AND diary_id IN (?) AND state = ?`, userId, diaryId, StateOnline).Find(&data)
	if len(data) > 0 {
		for _, v := range diaryId { // 对比 diaryId 是否存在数据表中
			for _, item := range data {
				if item.DiaryId == v {
					result[v] = true
				}
			}
			if _, ok := result[v]; !ok {
				result[v] = false
			}
		}
	} else {
		for _, v := range diaryId {
			result[v] = false
		}
	}
	return result
}
