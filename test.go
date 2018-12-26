package main

import (
	"git.chinawayltd.com/golib/tools/courier/enumeration"
	"time"
)

// 测试用表
//go:generate ast2template gen model Student --database DBEtrip --table-name student1
// @def primary ID
// @def unique_index I_name Name
// @def unique_index I_idcard IDCardNo
// @def index I_site Site Age
type Student struct {
	// 主键
	ID uint64 `gorm:"column:F_id"`
	// 姓名
	Name     string           `gorm:"column:F_name"`
	Pwd      string           `gorm:"column:F_pwd"`
	Site     string           `gorm:"column:F_site"`
	Age      int              `gorm:"column:F_age"`
	IDCardNo int              `gorm:"column:F_id_card_no"`
	Enabled  enumeration.Bool `gorm:"column:F_enabled"`
	// 创建时间
	CreateTime time.Time `gorm:"column:F_create_time"`
	// 更新时间
	UpdateTime time.Time `gorm:"column:F_update_time"`
}
