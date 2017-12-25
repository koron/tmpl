package tmpl

import "github.com/alecthomas/template"

var funcmap = template.FuncMap{}

func AddFunc(name string, fn interface{}) {
	funcmap[name] = fn
}
