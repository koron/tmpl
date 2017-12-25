package tmpl

import "text/template"

var funcmap = template.FuncMap{}

func AddFunc(name string, fn interface{}) {
	funcmap[name] = fn
}
