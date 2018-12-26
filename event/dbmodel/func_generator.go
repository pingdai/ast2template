package dbmodel

import (
	"fmt"
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"my_github/ast2template/codegen"
	"my_github/ast2template/codegen/processrx"
)

type Config struct {
	StructName string
	TableName  string
}

type DBFuncGenerator struct {
	Config
	pkgImportPath string
	program       *loader.Program
	model         *Model
}

func (g *DBFuncGenerator) Defaults() {
	if g.TableName == "" {
		g.TableName = codegen.ToLowerSnakeCase(g.StructName)
	}
}

func (g *DBFuncGenerator) Load(cwd string) {
	ldr := loader.Config{
		AllowErrors: true,
		ParserMode:  parser.ParseComments,
	}

	pkgImportPath := codegen.GetPackageImportPath(cwd)
	ldr.Import(pkgImportPath)

	p, err := ldr.Load()
	if err != nil {
		panic(err)
	}

	g.pkgImportPath = pkgImportPath
	g.program = p

	g.Defaults()
}

func (g *DBFuncGenerator) Process( /*ignoreTable bool*/ ) {
	for pkg, pkgInfo := range g.program.AllPackages {
		if pkg.Path() != g.pkgImportPath {
			continue
		}
		for ident, obj := range pkgInfo.Defs {
			if typeName, ok := obj.(*types.TypeName); ok {
				if typeName.Name() == g.StructName {
					if typeStruct, ok := typeName.Type().Underlying().(*types.Struct); ok {
						comments := processrx.CommentsOf(g.program.Fset, ident, pkgInfo.Files...)

						g.model = &Model{
							Pkg:            pkg,
							Name:           g.StructName,
							TableName:      g.TableName,
							UniqueIndex:    make(map[string][]Field),
							NormalIndex:    make(map[string][]Field),
							FuncMapContent: make(map[string]string),
						}
						g.model.collectIndexFromComments(comments)
						g.model.collectInfoFromStructType(typeStruct)

						for _, v := range g.model.PrimaryIndex {
							fmt.Printf("主键索引:%+v\n", v)
						}
						for key, fields := range g.model.NormalIndex {
							fmt.Printf("普通索引:%s\n", key)
							for _, v := range fields {
								fmt.Printf("	%+v\n", v)
							}
						}
						for key, fields := range g.model.UniqueIndex {
							fmt.Printf("唯一索引:%s\n", key)
							for _, v := range fields {
								fmt.Printf("	%+v\n", v)
							}
						}
						for _, v := range g.model.Fields {
							fmt.Printf("Field:	%+v\n", v)
						}

						g.model.Output(pkgInfo.Pkg.Name(), false)
					}
				}
			}
		}
	}
}

func (g *DBFuncGenerator) Output(cwd string) codegen.Outputs {
	return make(codegen.Outputs)
}
