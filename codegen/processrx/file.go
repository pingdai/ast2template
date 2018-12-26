package processrx

import "go/ast"

func FileOf(targetNode ast.Node, files ...*ast.File) *ast.File {
	for _, file := range files {
		if file.Pos() <= targetNode.Pos() && file.End() > targetNode.Pos() {
			return file
		}
	}
	return nil
}
