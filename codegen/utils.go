package codegen

import "go/build"

func GetPackageImportPath(dir string) string {
	pkg, err := build.ImportDir(dir, build.FindOnly)
	if err != nil {
		panic(err)
	}
	return pkg.ImportPath
}
