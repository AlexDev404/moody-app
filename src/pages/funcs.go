package pages

import (
	"html/template"
	"unicode"
	"unicode/utf8"
)

// Funcs are the template functions the app's templates use.
var Funcs = template.FuncMap{
	"CapitalizeFirst": CapitalizeFirst,
}

// CapitalizeFirst upper-cases the first letter of s and leaves the rest alone.
func CapitalizeFirst(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}
