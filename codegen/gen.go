package codegen

import "os"

type Generator interface {
	Load(cwd string)
	Process()
	Output()
}

func Generate(generator Generator) {
	cwd, _ := os.Getwd()
	generator.Load(cwd)
	generator.Process()
	generator.Output()
}
