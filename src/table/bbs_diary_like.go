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

// BbsDiaryLike 社区动态点赞
type BbsDiaryLike struct {
	Id         int        `gorm:"column:id;type:int(11);unsigned;not null;AUTO_INCREMENT"`              // 主键 id
	DiaryId    int        `gorm:"column:diary_id;type:int(11);unsigned;not null DEFAULT 0"`             // bbs_diary.id
	UserId     int        `gorm:"column:user_id;type:int(11);unsigned;not null DEFAULT 0"`              // 用户 id
	State      int        `gorm:"column:state;type:tinyint(1);not null DEFAULT 0"`                      // 状态, -1=下线, 1=上线
	CreateTime *time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"` // 创建时间
}

// TableName 获取数据表名
func (tb *BbsDiaryLike) TableName() string {
	return "bbs_diary_like"
}
