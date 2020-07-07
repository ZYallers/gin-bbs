package model

import (
	app "src/config"
	"src/table"
)

// bbsDiaryComment 社区动态评论, 将来可能会扩展到多个模块的评论
type bbsDiaryComment struct {
	baseModel
}

func NewBbsDiaryComment() *bbsDiaryComment {
	return &bbsDiaryComment{}
}

// InsertComment 添加评论
func (dc *bbsDiaryComment) InsertComment(comment table.BbsDiaryComment) int {
	var id int
	if dc.GetBbs().Create(&comment).Error == nil {
		id = comment.Id
	}
	return id
}

// FindComment 获取上线的评论列表
func (dc *bbsDiaryComment) FindComment(diaryId int, offset int, limit int) []table.BbsDiaryComment {
	var result []table.BbsDiaryComment
	dc.GetBbs().Select(`id,diary_id,user_id,content,create_time`).Where(`diary_id = ? AND state = ?`, diaryId, StateOnline).Limit(limit).Offset(offset).Order(`id DESC`).Find(&result)
	return result
}

// Template 评论模板
func (bbsDiaryComment) Template() []string {
	return app.Bbs.CommentTemplate
}

// GetCounter 获取上线的评论总数
func (dc *bbsDiaryComment) GetCounter(diaryId int) (row int) {
	dc.GetBbs().Model(&table.BbsDiaryComment{}).Where(`diary_id = ? AND state = ?`, diaryId, StateOnline).Count(&row)
	return row
}
