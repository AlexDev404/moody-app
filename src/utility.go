package main

import (
	"html/template"
	"strings"
	"unicode"
)

var funcMap = template.FuncMap{
	"CapitalizeFirst": CapitalizeFirst,
}

func CapitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	s = strings.ToLower(s[1:])
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
