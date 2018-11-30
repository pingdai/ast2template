package dbmodel

import "my_github/ast2template/codegen"

type Config struct {
	StructName string
	TableName  string
	Database   string
}

type DBFuncGenerator struct {
	Config
}

func (g *DBFuncGenerator) Load(cwd string) {

}

func (g *DBFuncGenerator) Process() {

}

func (g *DBFuncGenerator) Output(cwd string) codegen.Outputs {
	outputs := codegen.Outputs{}

	return outputs
}
