package codegen

import "strings"

func DeVendor(importPath string) string {
	parts := strings.Split(importPath, "/vendor/")
	return parts[len(parts)-1]
}
