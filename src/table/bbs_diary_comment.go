// Copyright (c) 2020 HXS R&D Technologies, Inc.
//
// @Author zhongyongbiao
// @Version 1.0.0
// @Time 2020/5/20 下午4:06
// @Software GoLand
package table

import (
	"time"
)

// BbsDiaryComment 社区动态评论
type BbsDiaryComment struct {
	Id         int        `gorm:"column:id;type:int(11);unsigned;not null;AUTO_INCREMENT"`              // 主键 id
	DiaryId    int        `gorm:"column:diary_id;type:int(11);unsigned;not null DEFAULT 0"`             // bbs_diary.id
	UserId     int        `gorm:"column:user_id;type:int(11);unsigned;not null DEFAULT 0"`              // 用户 id
	Content    string     `gorm:"column:content;type:text"`                                             // 内容
	State      int        `gorm:"column:state;type:tinyint(1);not null DEFAULT 0"`                      // 状态, -2=用户删除, -1=下线, 0=待审 1=审核通过
	CreateTime *time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"` // 创建时间
	UpdateTime *time.Time `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"` // 更新时间
}

// TableName 获取数据表名
func (tb *BbsDiaryComment) TableName() string {
	return "bbs_diary_comment"
}
