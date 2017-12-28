package tmpl

import (
	"strings"
	"unicode"

	"github.com/fatih/camelcase"
)

func init() {
	AddFunc("upperCamel", funcUpperCamel)
	AddFunc("lowerCamel", funcLowerCamel)
}

func camel(s string, fn func(rune) rune) string {
	list := camelcase.Split(s)
	if len(list) == 0 {
		return s
	}
	first := true
	list[0] = strings.Map(func(r rune) rune {
		if first {
			first = false
			return fn(r)
		}
		return r
	}, list[0])
	return strings.Join(list, "")
}

func funcUpperCamel(s string) string {
	return camel(s, unicode.ToUpper)
}

func funcLowerCamel(s string) string {
	return camel(s, unicode.ToLower)
}
