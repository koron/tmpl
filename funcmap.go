package tmpl

import "text/template"

var funcmap = template.FuncMap{}

// AddFunc add a function which can be used from template.
func AddFunc(name string, fn interface{}) {
	funcmap[name] = fn
}
