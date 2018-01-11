package tmpl

import "github.com/koron/tmpl/internal/gosrc"

// SourceGosrc is a source for LoadGosrc(). Default is "." (current directory)
var SourceGosrc string

// LoadGosrc loads a file or a directory as go source code.
func LoadGosrc() (interface{}, error) {
	path := SourceGosrc
	if path == "" {
		path = "."
	}
	pkg, err := gosrc.Read(path)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}
