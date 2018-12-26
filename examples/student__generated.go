package examples

import (
	"time"

	"github.com/jinzhu/gorm"
)

type StudentList []Student

func (s Student) TableName() string {
	table_name := "t_student"
	return table_name
}

func (sl *StudentList) BatchFetchByClassList(db *gorm.DB, classList []int) error {
	if len(classList) == 0 {
		return nil
	}

	err := db.Table(Student{}.TableName()).Where("F_class in (?) and F_enabled = ?", classList, 1).Find(sl).Error

	return err
}

func (sl *StudentList) BatchFetchByIDList(db *gorm.DB, iDList []uint64) error {
	if len(iDList) == 0 {
		return nil
	}

	err := db.Table(Student{}.TableName()).Where("F_id in (?) and F_enabled = ?", iDList, 1).Find(sl).Error

	return err
}

func (sl *StudentList) BatchFetchByNameList(db *gorm.DB, nameList []string) error {
	if len(nameList) == 0 {
		return nil
	}

	err := db.Table(Student{}.TableName()).Where("F_name in (?) and F_enabled = ?", nameList, 1).Find(sl).Error

	return err
}

func (sl *StudentList) BatchFetchByNoList(db *gorm.DB, noList []string) error {
	if len(noList) == 0 {
		return nil
	}

	err := db.Table(Student{}.TableName()).Where("F_no in (?) and F_enabled = ?", noList, 1).Find(sl).Error

	return err
}

func (s *Student) Create(db *gorm.DB) error {

	if s.CreateTime.IsZero() {
		s.CreateTime = time.Now()
	}

	if s.UpdateTime.IsZero() {
		s.UpdateTime = time.Now()
	}

	s.Enabled = int(1)
	err := db.Table(s.TableName()).Create(s).Error

	return err
}

func (sl *StudentList) FetchByClass(db *gorm.DB, class int) error {
	err := db.Table(Student{}.TableName()).Where("F_class = ? and F_enabled = ?", class, 1).Find(sl).Error

	return err
}

func (s *Student) FetchByID(db *gorm.DB) error {
	err := db.Table(s.TableName()).Where("F_id = ? and F_enabled = ?", s.ID, 1).Find(s).Error

	return err
}

func (s *Student) FetchByIDForUpdate(db *gorm.DB) error {
	err := db.Table(s.TableName()).Where("F_id = ? and F_enabled = ?", s.ID, 1).Set("gorm:query_option", "FOR UPDATE").Find(s).Error

	return err
}

func (sl *StudentList) FetchByName(db *gorm.DB, name string) error {
	err := db.Table(Student{}.TableName()).Where("F_name = ? and F_enabled = ?", name, 1).Find(sl).Error

	return err
}

func (s *Student) FetchByNameAndAddress(db *gorm.DB) error {
	err := db.Table(s.TableName()).Where("F_name = ? and F_address = ? and F_enabled = ?", s.Name, s.Address, 1).Find(s).Error

	return err
}

func (s *Student) FetchByNameAndAddressForUpdate(db *gorm.DB) error {
	err := db.Table(s.TableName()).Where("F_name = ? and F_address = ? and F_enabled = ?", s.Name, s.Address, 1).Set("gorm:query_option", "FOR UPDATE").Find(s).Error

	return err
}

func (s *Student) FetchByNo(db *gorm.DB) error {
	err := db.Table(s.TableName()).Where("F_no = ? and F_enabled = ?", s.No, 1).Find(s).Error

	return err
}

func (s *Student) FetchByNoForUpdate(db *gorm.DB) error {
	err := db.Table(s.TableName()).Where("F_no = ? and F_enabled = ?", s.No, 1).Set("gorm:query_option", "FOR UPDATE").Find(s).Error

	return err
}

func (sl *StudentList) FetchList(db *gorm.DB, size, offset int32, query ...map[string]interface{}) (int32, error) {
	var count int32
	if len(query) == 0 {
		query = append(query, map[string]interface{}{"F_enabled": 1})
	} else {
		if _, ok := query[0]["F_enabled"]; !ok {
			query[0]["F_enabled"] = 1
		}
	}

	if size <= 0 {
		size = -1
		offset = -1
	}
	var err error

	err = db.Table(Student{}.TableName()).Where(query[0]).Count(&count).Limit(size).Offset(offset).Order("F_create_time desc").Find(sl).Error

	return int32(count), err
}

func (s *Student) SoftDeleteByID(db *gorm.DB) error {
	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = 2

	if s.UpdateTime.IsZero() {
		s.UpdateTime = time.Now()
	}

	err := db.Table(s.TableName()).Where("F_id = ? and F_enabled = ?", s.ID, 1).Updates(updateMap).Error

	return err
}

func (s *Student) SoftDeleteByNameAndAddress(db *gorm.DB) error {
	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = 2

	if s.UpdateTime.IsZero() {
		s.UpdateTime = time.Now()
	}

	err := db.Table(s.TableName()).Where("F_name = ? and F_address = ? and F_enabled = ?", s.Name, s.Address, 1).Updates(updateMap).Error

	return err
}

func (s *Student) SoftDeleteByNo(db *gorm.DB) error {
	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = 2

	if s.UpdateTime.IsZero() {
		s.UpdateTime = time.Now()
	}

	err := db.Table(s.TableName()).Where("F_no = ? and F_enabled = ?", s.No, 1).Updates(updateMap).Error

	return err
}

func (s *Student) UpdateByIDWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	if _, ok := updateMap["F_update_time"]; !ok {
		updateMap["F_update_time"] = time.Now()

	}
	dbRet := db.Table(s.TableName()).Where("F_id = ? and F_enabled = ?", s.ID, 1).Updates(updateMap)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			findErr := db.Table(s.TableName()).Where("F_id = ? and F_enabled = ?", s.ID, 1).Find(&Student{}).Error
			if findErr != nil {
				return findErr
			}
			//存在有效数据记录，返回成功
			return nil
		} else {
			return nil
		}
	}
}

func (s *Student) UpdateByIDWithStruct(db *gorm.DB) error {

	if s.UpdateTime.IsZero() {
		s.UpdateTime = time.Now()
	}

	dbRet := db.Table(s.TableName()).Where("F_id = ? and F_enabled = ?", s.ID, 1).Updates(s)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			findErr := db.Table(s.TableName()).Where("F_id = ? and F_enabled = ?", s.ID, 1).Find(&Student{}).Error
			if findErr != nil {
				return findErr
			}
			//存在有效数据记录，返回成功
			return nil
		} else {
			return nil
		}
	}
}

func (s *Student) UpdateByNameAndAddressWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	if _, ok := updateMap["F_update_time"]; !ok {
		updateMap["F_update_time"] = time.Now()

	}
	dbRet := db.Table(s.TableName()).Where("F_name = ? and F_address = ? and F_enabled = ?", s.Name, s.Address, 1).Updates(updateMap)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			findErr := db.Table(s.TableName()).Where("F_name = ? and F_address = ? and F_enabled = ?", s.Name, s.Address, 1).Find(&Student{}).Error
			if findErr != nil {
				return findErr
			}
			//存在有效数据记录，返回成功
			return nil
		} else {
			return nil
		}
	}
}

func (s *Student) UpdateByNameAndAddressWithStruct(db *gorm.DB) error {

	if s.UpdateTime.IsZero() {
		s.UpdateTime = time.Now()
	}

	dbRet := db.Table(s.TableName()).Where("F_name = ? and F_address = ? and F_enabled = ?", s.Name, s.Address, 1).Updates(s)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			findErr := db.Table(s.TableName()).Where("F_name = ? and F_address = ? and F_enabled = ?", s.Name, s.Address, 1).Find(&Student{}).Error
			if findErr != nil {
				return findErr
			}
			//存在有效数据记录，返回成功
			return nil
		} else {
			return nil
		}
	}
}

func (s *Student) UpdateByNoWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	if _, ok := updateMap["F_update_time"]; !ok {
		updateMap["F_update_time"] = time.Now()

	}
	dbRet := db.Table(s.TableName()).Where("F_no = ? and F_enabled = ?", s.No, 1).Updates(updateMap)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			findErr := db.Table(s.TableName()).Where("F_no = ? and F_enabled = ?", s.No, 1).Find(&Student{}).Error
			if findErr != nil {
				return findErr
			}
			//存在有效数据记录，返回成功
			return nil
		} else {
			return nil
		}
	}
}

func (s *Student) UpdateByNoWithStruct(db *gorm.DB) error {

	if s.UpdateTime.IsZero() {
		s.UpdateTime = time.Now()
	}

	dbRet := db.Table(s.TableName()).Where("F_no = ? and F_enabled = ?", s.No, 1).Updates(s)
	err := dbRet.Error
	if err != nil {
		return err
	} else {
		if dbRet.RowsAffected == 0 {
			findErr := db.Table(s.TableName()).Where("F_no = ? and F_enabled = ?", s.No, 1).Find(&Student{}).Error
			if findErr != nil {
				return findErr
			}
			//存在有效数据记录，返回成功
			return nil
		} else {
			return nil
		}
	}
}
