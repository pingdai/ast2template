package dbmodel

import (
	"regexp"
	"strings"
)

var (
	defRegexp = regexp.MustCompile(`@def ([^\n]+)`)
)

func defSplit(def string) (defs []string) {
	vs := strings.Split(def, " ")
	for _, s := range vs {
		if s != "" {
			defs = append(defs, s)
		}
	}
	return
}

func defToField(defs []string) []Field {
	fields := make([]Field, 0, len(defs))
	for _, v := range defs {
		fields = append(fields, Field{Name: v})
	}

	return fields
}
