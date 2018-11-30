package dbmodel

import (
	"go/parser"
	"golang.org/x/tools/go/loader"
	"my_github/ast2template/codegen"
)

type Config struct {
	StructName string
	TableName  string
	Database   string
}

type DBFuncGenerator struct {
	Config
	pkgImportPath string
	program       *loader.Program
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

func (g *DBFuncGenerator) Process() {

}

func (g *DBFuncGenerator) Output(cwd string) codegen.Outputs {
	outputs := codegen.Outputs{}

	return outputs
}
