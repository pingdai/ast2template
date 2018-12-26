package examples

import (
	"time"
)

// 测试用表
//go:generate ast2template gen model Student --table-name student
// @def primary ID
// @def unique_index U_no No
// @def unique_index U_name_address Name Address
// @def index I_site Class
type Student struct {
	// 主键
	ID uint64 `gorm:"column:F_id"`
	// 姓名
	Name string `gorm:"column:F_name"`
	// 学号
	No string `gorm:"column:F_no"`
	// 班级
	Class int `gorm:"column:F_class"`
	// 年龄
	Age int `gorm:"column:F_age"`
	// 住址
	Address string `gorm:"column:F_address"`
	// 是否有效，1是，2否
	Enabled int `gorm:"column:F_enabled"`
	// 更新时间
	UpdateTime time.Time `gorm:"column:F_update_time"`
	// 创建时间
	CreateTime time.Time `gorm:"column:F_create_time"`
}
