// Copyright (c) 2020 HXS R&D Technologies, Inc.
//
// @Author zhongyongbiao
// @Version 1.0.0
// @Time 2020/5/20 下午4:06
// @Software GoLand
package table

type BbsDiary struct {
	Id          int    `gorm:"column:id;type:int(11);unsigned;not null;primary_key;auto_increment" json:"id"`                                    // 主键
	UserId      int    `gorm:"column:user_id;type:int(11);unsigned;not null;default:'0';index:idx_userid" json:"user_id"`                        // 用户id
	Content     string `gorm:"column:content;type:varchar(500);not null;default:''" json:"content"`                                              // 内容
	SubData     string `gorm:"column:sub_data;type:varchar(400);not null;default:''" json:"sub_data"`                                            // 子类型 1-图片集，2-视频
	SubType     int    `gorm:"column:sub_type;type:tinyint(1);unsigned;not null;default:'1'" json:"sub_type"`                                    // 用户id
	DataId      int    `gorm:"column:data_id;type:int(11);unsigned;not null;default:'0'" json:"data_id"`                                         // 内容Id (例如 圈子Id)
	DataType    int    `gorm:"column:data_type;type:tinyint(1);unsigned;not null;default:'1'" json:"data_type"`                                  // 类型 1-圈子
	State       int    `gorm:"column:state;type:tinyint(1);unsigned;not null;default:'0'" json:"state"`                                          // 状态 -2-用户删除，-1-下线，0-代审，1-审核通过
	AdminUserId int    `gorm:"column:admin_user_id;type:int(11);unsigned;not null;default:'0'" json:"admin_user_id"`                             // 后台管理员id 审核
	CreateTime  string `gorm:"column:create_time;type:timestamp;not null;default:'1999-01-01 00:00:00';index:idx_createtime" json:"create_time"` // 创建时间
	UpdateTime  string `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"update_time"`                          // 更新时间
}

func (tb *BbsDiary) TableName() string {
	return "bbs_diary"
}
