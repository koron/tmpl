package tmpl

import (
	"errors"
	"io"
	"path/filepath"
	"text/template"
)

// DataFunc provides data for tmpl.Execute()
type DataFunc func() (interface{}, error)

// Execute executes templates set.
//
// If dataFunc is omitted, YAML is loaded by tmpl.LoadYAML
func Execute(dataFunc DataFunc, out io.Writer, tmplFiles ...string) error {
	if len(tmplFiles) == 0 {
		return errors.New("no template files")
	}
	name := filepath.Base(tmplFiles[0])

	tmpl, err := template.New(name).Funcs(funcmap).ParseFiles(tmplFiles...)
	if err != nil {
		return err
	}

	if dataFunc == nil {
		dataFunc = LoadYAML
	}
	data, err := dataFunc()
	if err != nil {
		return err
	}

	err = tmpl.Execute(out, data)
	if err != nil {
		return err
	}
	return nil
}
