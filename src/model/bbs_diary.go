package model

import (
	"src/abs"
	"src/library/tool"
	"src/table"
	"strconv"
)

type bbsDiary struct {
	abs.Model
}

type fixDiary struct {
	table.BbsDiary
	CircleName string
	Icon       string
}

func NewBbsDiary() *bbsDiary {
	return &bbsDiary{}
}

func (bd *bbsDiary) InsertDiary(diary table.BbsDiary) int {
	if bd.GetBbs().Create(&diary).Error == nil {
		return diary.Id
	}
	return 0
}

func (bd *bbsDiary) DeleteDiary(id int) int64 {
	return bd.GetBbs().Model(table.BbsDiary{Id: id}).Updates(table.BbsDiary{State: -2}).RowsAffected
}

func (bd *bbsDiary) FindDiary(ids []string, fields string) []fixDiary {
	var diarySlice []fixDiary
	bd.GetBbs().Table("bbs_diary").Where("bbs_diary.id in (?)", ids).
		Select("bbs_diary.*, bbs_circle.title as circle_name, bbs_circle.icon").
		Joins("left join bbs_circle on bbs_circle.id = bbs_diary.data_id").Find(&diarySlice)
	return diarySlice
}

func (bd *bbsDiary) FindCircleDiary(circleId string) []string {
	var idSlice []string
	diarySlice := []struct{ Id int }{}
	bd.GetBbs().Table("bbs_diary").Where("data_type=1 and data_id=? and state=1", circleId).Select("id").Find(&diarySlice)
	for _, v := range diarySlice {
		idSlice = append(idSlice, tool.FixPad(strconv.Itoa(v.Id)))
	}
	return idSlice
}
