package dbmodel

import (
	"fmt"
	"strings"
)

func ParseTagSetting(str string) map[string][]string {
	tags := strings.Split(str, ";")
	setting := map[string][]string{}
	for _, value := range tags {
		v := strings.Split(value, ":")
		k := strings.TrimSpace(strings.ToUpper(v[0]))
		if _, ok := setting[k]; !ok {
			setting[k] = make([]string, 0, 10)
		}
		if len(v) == 2 {
			setting[k] = append(setting[k], v[1])
		} else {
			setting[k] = append(setting[k], k)
		}
	}
	return setting
}

func replaceUpperWithLowerAndUnderscore(src string) string {
	var dst string
	for index, letter := range src {
		if index == 0 && isUpperLetter(letter) {
			dst += fmt.Sprintf("%c", letter)
		} else if isUpperLetter(letter) {
			dst += fmt.Sprintf("_%c", letter)
		} else {
			dst += fmt.Sprintf("%c", letter)
		}
	}

	return strings.ToLower(dst)
}

// fetchBaseInfoOfGenFuncForNormalIndex fetch part of function name, function input param,
// orm query format, orm query parameters.
func fetchBaseInfoOfGenFuncForNormalIndex(indexList []Field) *BaseInfoOfGenCode {
	var partFuncName, inputParam, ormQueryFormat, ormQueryParam string
	for _, field := range indexList {
		if len(partFuncName) == 0 {
			partFuncName = field.Name
			inputParam = fmt.Sprintf("db *gorm.DB, %s %s", convertFirstLetterToLower(field.Name), field.Type)
			ormQueryFormat = fmt.Sprintf("%s = ?", field.DbFieldName)
			ormQueryParam = fmt.Sprintf("%s", convertFirstLetterToLower(field.Name))
		} else {
			partFuncName += "And" + field.Name
			inputParam += fmt.Sprintf(", %s %s", convertFirstLetterToLower(field.Name), field.Type)
			ormQueryFormat += fmt.Sprintf(" and %s = ?", field.DbFieldName)
			ormQueryParam += fmt.Sprintf(", %s", convertFirstLetterToLower(field.Name))
		}
	}

	return &BaseInfoOfGenCode{
		PartFuncName:   partFuncName,
		FuncInputParam: inputParam,
		OrmQueryFormat: ormQueryFormat,
		OrmQueryParam:  ormQueryParam,
	}
}

func fetchBaseInfoOfGenFuncForUniqueIndex(model *Model, indexList []Field) *BaseInfoOfGenCode {
	var partFuncName, inputParam, ormQueryFormat, ormQueryParam string
	inputParam = fmt.Sprintf("%s", "db *gorm.DB")
	for _, field := range indexList {
		if len(partFuncName) == 0 {
			partFuncName = field.Name
			ormQueryFormat = fmt.Sprintf("%s = ?", field.DbFieldName)
			ormQueryParam = fmt.Sprintf("%s.%s", fetchUpLetter(model.Name), field.Name)
		} else {
			partFuncName += "And" + field.Name
			ormQueryFormat += fmt.Sprintf(" and %s = ?", field.DbFieldName)
			ormQueryParam += fmt.Sprintf(", %s.%s", fetchUpLetter(model.Name), field.Name)
		}
	}

	return &BaseInfoOfGenCode{
		PartFuncName:   partFuncName,
		FuncInputParam: inputParam,
		OrmQueryFormat: ormQueryFormat,
		OrmQueryParam:  ormQueryParam,
	}
}
