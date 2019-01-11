package enum

import "golang.org/x/tools/go/loader"

type EnumGenerator struct {
	Filters       []string
	pkgImportPath string
	program       *loader.Program
	enumScanner   *processrx.EnumScanner
}
