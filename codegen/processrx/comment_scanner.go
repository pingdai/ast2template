package processrx

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

func CommentsOf(fileSet *token.FileSet, targetNode ast.Node, files ...*ast.File) string {
	file := FileOf(targetNode, files...)
	if file == nil {
		return ""
	}
	commentScanner := NewCommentScanner(fileSet, file)
	doc := commentScanner.CommentsOf(targetNode)
	if doc != "" {
		return doc
	}
	return doc
}

func NewCommentScanner(fileSet *token.FileSet, file *ast.File) *CommentScanner {
	commentMap := ast.NewCommentMap(fileSet, file, file.Comments)

	return &CommentScanner{
		file:       file,
		CommentMap: commentMap,
	}
}

type CommentScanner struct {
	file       *ast.File
	CommentMap ast.CommentMap
}

func (scanner *CommentScanner) CommentsOf(targetNode ast.Node) string {
	commentGroupList := scanner.CommentGroupListOf(targetNode)
	return StringifyCommentGroup(commentGroupList...)
}

func (scanner *CommentScanner) CommentGroupListOf(targetNode ast.Node) (commentGroupList []*ast.CommentGroup) {
	if targetNode == nil {
		return
	}

	switch targetNode.(type) {
	case *ast.File, *ast.Field, ast.Stmt, ast.Decl:
		if comments, ok := scanner.CommentMap[targetNode]; ok {
			commentGroupList = comments
		}
	case ast.Spec:
		// Spec should merge with comments of its parent gen decl when empty
		if comments, ok := scanner.CommentMap[targetNode]; ok {
			commentGroupList = append(commentGroupList, comments...)
		}

		if len(commentGroupList) == 0 {
			for node, comments := range scanner.CommentMap {
				if genDecl, ok := node.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if targetNode == spec {
							commentGroupList = append(commentGroupList, comments...)
						}
					}
				}
			}
		}
	default:
		// find nearest parent node which have comments
		{
			var deltaPos token.Pos
			var parentNode ast.Node

			deltaPos = -1

			ast.Inspect(scanner.file, func(node ast.Node) bool {
				switch node.(type) {
				case *ast.Field, ast.Decl, ast.Spec, ast.Stmt:
					if targetNode.Pos() >= node.Pos() && targetNode.End() <= node.End() {
						nextDelta := targetNode.Pos() - node.Pos()
						if deltaPos == -1 || (nextDelta <= deltaPos) {
							deltaPos = nextDelta
							parentNode = node
						}
					}
				}
				return true
			})

			if parentNode != nil {
				commentGroupList = scanner.CommentGroupListOf(parentNode)
			}
		}
	}

	sort.Sort(ByCommentPos(commentGroupList))
	return
}

type ByCommentPos []*ast.CommentGroup

func (a ByCommentPos) Len() int {
	return len(a)
}

func (a ByCommentPos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByCommentPos) Less(i, j int) bool {
	return a[i].Pos() < a[j].Pos()
}

func StringifyCommentGroup(commentGroupList ...*ast.CommentGroup) (comments string) {
	if len(commentGroupList) == 0 {
		return ""
	}
	for _, commentGroup := range commentGroupList {
		for _, line := range strings.Split(commentGroup.Text(), "\n") {
			if strings.HasPrefix(line, "go:generate") {
				continue
			}
			comments = comments + "\n" + line
		}
	}
	return strings.TrimSpace(comments)
}
