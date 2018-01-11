package tmpl

import (
	"errors"
	"io"
	"path/filepath"
	"text/template"
)

type LoadDataFunc func() (interface{}, error)

// Execute executes templates set.
//
// If loadData is omitted, YAML is loaded by tmpl.LoadYAML
func Execute(loadData LoadDataFunc, out io.Writer, tmplFiles ...string) error {
	if len(tmplFiles) == 0 {
		return errors.New("no template files")
	}
	name := filepath.Base(tmplFiles[0])

	tmpl, err := template.New(name).Funcs(funcmap).ParseFiles(tmplFiles...)
	if err != nil {
		return err
	}

	if loadData == nil {
		loadData = LoadYAML
	}
	data, err := loadData()
	if err != nil {
		return err
	}

	err = tmpl.Execute(out, data)
	if err != nil {
		return err
	}
	return nil
}
