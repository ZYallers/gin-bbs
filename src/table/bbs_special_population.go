package table

type BbsSpecialPopulation struct {
	Id          int    `gorm:"column:id;type:int(11);unsigned;not null;primary_key;auto_increment" json:"id"`                     // 主键
	UserId      int    `gorm:"column:user_id;type:int(11);unsigned;not null;default:'0';unique_index:uniq_userid" json:"user_id"` // 主用户Id
	CircleData  string `gorm:"column:circle_data;type:varchar(100);not null;default:''" json:"circle_data"`                       // 优先可见圈子
	State       int    `gorm:"column:state;type:tinyint(1);unsigned;not null;default:'1'" json:"state"`                           // 状态 -1-删除，1-上线
	AdminUserId int    `gorm:"column:admin_user_id;type:int(11);unsigned;not null;default:'0'" json:"admin_user_id"`              // 后台管理员id
	CreateTime  string `gorm:"column:create_time;type:timestamp;not null;default:'1999-01-01 00:00:00'" json:"create_time"`       // 创建时间
	UpdateTime  string `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"update_time"`           // 更新时间
}

func (tb *BbsSpecialPopulation) TableName() string {
	return "bbs_special_population"
}
