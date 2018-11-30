package codegen

import "os"

type Outputs map[string]string

func (outputs Outputs) WriteFiles() {
	//for filename, content := range outputs {
	//	outputs.WriteFile(filename, content)
	//}
}

type Generator interface {
	Load(cwd string)
	Process()
	Output(cwd string) Outputs
}

func Generate(generator Generator) {
	cwd, _ := os.Getwd()
	generator.Load(cwd)
	generator.Process()
	outputs := generator.Output(cwd)
	outputs.WriteFiles()
}
