package dbmodel

import (
	"reflect"
	"testing"
)

func TestParseTagSetting(t *testing.T) {
	str := `gorm:"column:F_name" sql:"type:bigint(64);not null;index:I_update_time"`
	tag := reflect.StructTag(str)
	rt := ParseTagSetting(tag.Get("gorm"))
	t.Logf("%+v", rt)
}
