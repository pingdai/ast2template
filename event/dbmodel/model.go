package dbmodel

import (
	"ast2template/codegen"
	"ast2template/codegen/processrx"
	"fmt"
	"go/build"
	"go/types"
	"golang.org/x/tools/imports"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

type Field struct {
	Name         string
	Type         string
	DbFieldName  string
	IsEnable     bool
	IsCreateTime bool
}

type Model struct {
	Pkg                 *types.Package
	TableName           string
	Name                string
	Fields              []*Field
	IsDbModel           bool
	UniqueIndex         map[string][]Field
	NormalIndex         map[string][]Field
	PrimaryIndex        []Field
	HasCreateTimeField  bool
	HasUpdateTimeField  bool
	HasEnabledField     bool
	EnabledFieldType    string
	CreateTimeFieldType string
	UpdateTimeFieldType string
	DbCreateTimeField   string
	DbEnabledField      string
	DbUpdateTimeField   string
	FuncMapContent      map[string]string
	Deps                []string
}

func (model *Model) collectIndexFromComments(comments string) {
	matches := defRegexp.FindAllStringSubmatch(comments, -1)

	for _, subMatch := range matches {
		if len(subMatch) == 2 {
			defs := defSplit(subMatch[1])

			switch strings.ToLower(defs[0]) {
			case "primary":
				if len(defs) < 2 {
					panic(fmt.Errorf("primary at lease 1 Field"))
				}

				model.PrimaryIndex = append(model.PrimaryIndex, defToField(defs[1:])...)
			case "index":
				if len(defs) < 3 {
					panic(fmt.Errorf("index at lease 1 Field"))
				}
				if model.NormalIndex == nil {
					model.NormalIndex = make(map[string][]Field)
				}
				if _, ok := model.NormalIndex[defs[1]]; ok {
					panic(fmt.Errorf("repeat index[%s]", defs[1]))
				}

				model.NormalIndex[defs[1]] = append(model.NormalIndex[defs[1]], defToField(defs[2:])...)
			case "unique_index":
				if len(defs) < 3 {
					panic(fmt.Errorf("unique Indexes at lease 1 Field"))
				}
				if model.UniqueIndex == nil {
					model.UniqueIndex = make(map[string][]Field)
				}
				if _, ok := model.UniqueIndex[defs[1]]; ok {
					panic(fmt.Errorf("repeat unique_index[%s]", defs[1]))
				}

				model.UniqueIndex[defs[1]] = append(model.UniqueIndex[defs[1]], defToField(defs[2:])...)
			}
		}
	}
}

func (model *Model) collectInfoFromStructType(typeStruct *types.Struct) {
	for i := 0; i < typeStruct.NumFields(); i++ {
		field := typeStruct.Field(i)
		tag := reflect.StructTag(typeStruct.Tag(i))

		gormSettings := ParseTagSetting(tag.Get("gorm"))
		if len(gormSettings) != 0 {
			if _, ok := gormSettings["-"]; ok {
				continue
			}
			model.IsDbModel = true
		} else {
			continue
		}

		var tmpField = Field{}
		if dbFieldName, ok := gormSettings["COLUMN"]; ok {
			tmpField.DbFieldName = dbFieldName[0]
		} else {
			tmpField.DbFieldName = tmpField.Name
		}
		tmpField.Name = field.Name()

		pkgPath, method := processrx.GetPkgImportPathAndExpose(field.Type().String())
		if pkgPath != "" {
			pkg, _ := build.Import(pkgPath, "", build.ImportComment)
			tmpField.Type = fmt.Sprintf("%s.%s", pkg.Name, method)
		} else {
			tmpField.Type = method
		}

		if pkgPath != "" {
			model.Deps = append(model.Deps, pkgPath)
		}

		model.Fields = append(model.Fields, &tmpField)
		if tmpField.Name == "Enabled" || tmpField.DbFieldName == "F_enabled" {
			tmpField.IsEnable = true
			model.HasEnabledField = true
			model.EnabledFieldType = tmpField.Type
			model.DbEnabledField = tmpField.DbFieldName
		} else if tmpField.Name == "UpdateTime" || tmpField.DbFieldName == "F_update_time" {
			model.HasUpdateTimeField = true
			model.UpdateTimeFieldType = tmpField.Type
			model.DbUpdateTimeField = tmpField.DbFieldName
		} else if tmpField.Name == "CreateTime" || tmpField.DbFieldName == "F_create_time" {
			tmpField.IsCreateTime = true
			model.HasCreateTimeField = true
			model.CreateTimeFieldType = tmpField.Type
			model.DbCreateTimeField = tmpField.DbFieldName
		}
	}

	// 校验索引中的值在结构体字段中是否存在
	for i, primary := range model.PrimaryIndex {
		var flag bool
		for _, field := range model.Fields {
			if primary.Name == field.Name {
				model.PrimaryIndex[i] = (*field)
				flag = true
				break
			}
		}
		if !flag {
			panic(fmt.Errorf("def field name[%s] not exist", primary.Name))
		}
	}
	for key, normalFields := range model.NormalIndex {
		for i, normal := range normalFields {
			var flag bool
			for _, field := range model.Fields {
				if normal.Name == field.Name {
					model.NormalIndex[key][i] = (*field)
					flag = true
					break
				}
			}
			if !flag {
				panic(fmt.Errorf("def field name[%s] not exist", normal.Name))
			}
		}
	}
	for key, uniqueFields := range model.UniqueIndex {
		for i, unique := range uniqueFields {
			var flag bool
			for _, field := range model.Fields {
				if unique.Name == field.Name {
					model.UniqueIndex[key][i] = (*field)
					flag = true
					break
				}
			}
			if !flag {
				panic(fmt.Errorf("def field name[%s] not exist", unique.Name))
			}
		}
	}
}

func (model *Model) Output(ignoreCreateTableNameFunc bool) {
	if err := genTableNameFunc(model, ignoreCreateTableNameFunc); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genCreateFunc(model); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	model.genCodeByNormalIndex()
	model.genCodeByUniqueIndex()
	model.genCodeByPrimaryKeyIndex()

	if err := genFetchListFunc(model); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	model.GenerateFile()
}

func (model *Model) GenerateFile() {
	var funcNameList []string

	// first part of file
	var tableName = "tableName"
	for key := range model.FuncMapContent {
		if key == tableName {
			continue
		}
		funcNameList = append(funcNameList, key)
	}

	contents := []string{
		model.FuncMapContent[tableName],
	}

	sort.Strings(funcNameList)

	for _, funcName := range funcNameList {
		contents = append(contents, model.FuncMapContent[funcName])
	}

	p, _ := build.Import(model.Pkg.Path(), "", build.FindOnly)
	cwd, _ := os.Getwd()
	path, _ := filepath.Rel(cwd, p.Dir)

	filename := path + "/" + replaceUpperWithLowerAndUnderscore(model.Name) + ".go"
	content := strings.Join(contents, "\n\n")
	bytes, err := imports.Process(filename, []byte(content), nil)
	if err != nil {
		panic(err)
	} else {
		content = string(bytes)
	}
	codegen.WriteFile(codegen.GeneratedSuffix(filename), content)
}

func (model *Model) genCodeByUniqueIndex() {
	for _, fieldList := range model.UniqueIndex {
		model.handleGenCodeForUniqueIndex(fieldList)
	}
}

func (model *Model) genCodeByPrimaryKeyIndex() {
	if len(model.PrimaryIndex) > 0 {
		model.handleGenCodeForUniqueIndex(model.PrimaryIndex)
	}
}

func (model *Model) genCodeByNormalIndex() {
	for _, fieldList := range model.NormalIndex {
		baseInfoGenCode := fetchBaseInfoOfGenFuncForNormalIndex(fieldList)
		if err := genFetchFuncByNormalIndex(model, baseInfoGenCode); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		model.handleGenFetchCodeBySubIndex(fieldList)
	}
}

func (model *Model) genBatchFetchFuncBySingleIndex(field Field) {
	if err := genBatchFetchFunc(model, field.Name, field.DbFieldName, field.Type); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}

func (model *Model) handleGenFetchCodeBySubIndex(fieldList []Field) {
	if len(fieldList) == 1 {
		model.genBatchFetchFuncBySingleIndex(fieldList[0])
	} else if len(fieldList) > 1 {
		// [x, y, z, e] Split to [x, y, z], [x, y], [x]
		for i := 1; i < len(fieldList); i++ {
			subSortFieldSlice := fieldList[:len(fieldList)-i]
			baseInfoGenCode := fetchBaseInfoOfGenFuncForNormalIndex(subSortFieldSlice)
			if err := genFetchFuncByNormalIndex(model, baseInfoGenCode); err != nil {
				fmt.Printf("%s\n", err.Error())
				os.Exit(1)
			}

			if len(subSortFieldSlice) == 1 {
				model.genBatchFetchFuncBySingleIndex(subSortFieldSlice[0])
			}

		}
	}
}

func (model *Model) handleGenCodeForUniqueIndex(sortFieldList []Field) {
	baseInfoGenCode := fetchBaseInfoOfGenFuncForUniqueIndex(model, sortFieldList)
	if err := genFetchFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	model.handleGenFetchCodeBySubIndex(sortFieldList)

	if err := genFetchForUpdateFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genUpdateWithStructFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genUpdateWithMapFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genSoftDeleteFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	// 物理删除
	//if err := genPhysicsDeleteFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
	//	fmt.Printf("%s\n", err.Error())
	//	os.Exit(1)
	//}
}
