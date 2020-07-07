// Copyright (c) 2020 HXS R&D Technologies, Inc.
//
// @Author zhongyongbiao
// @Version 1.0.0
// @Time 2020/5/20 下午4:06
// @Software GoLand
package table

type BbsCircle struct {
	Id           int    `gorm:"column:id;type:int(11);unsigned;not null;primary_key;auto_increment" json:"id"`               // 主键
	UserId       int    `gorm:"column:user_id;type:int(11);unsigned;not null;default:'0';index:idx_userid" json:"user_id"`   // 主用户Id
	Title        string `gorm:"column:title;type:varchar(50);not null;default:''" json:"title"`                              // 标题
	Introduction string `gorm:"column:introduction;type:varchar(60);not null;default:''" json:"introduction"`                // 简介
	Icon         string `gorm:"column:icon;type:varchar(100);not null;default:''" json:"icon"`                               // Icon
	Image        string `gorm:"column:image;type:varchar(100);not null;default:''" json:"image"`                             // 背景图
	Sort         int    `gorm:"column:sort;type:int(11);unsigned;not null;default:'0'" json:"sort"`                          // 排序
	State        int    `gorm:"column:state;type:tinyint(1);unsigned;not null;default:'-1'" json:"state"`                    // 状态 -1-下线，1-审核通过
	VisibleUser  int    `gorm:"column:visible_user;type:char(12);not null;default:'0'" json:"visible_user"`                  // 可见用户 位运算
	CreateTime   string `gorm:"column:create_time;type:timestamp;not null;default:'1999-01-01 00:00:00'" json:"create_time"` // 创建时间
	UpdateTime   string `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"update_time"`     // 更新时间
}

func (tb *BbsCircle) TableName() string {
	return "bbs_circle"
}
