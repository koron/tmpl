package tmpl

import (
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

func readYAML(r io.Reader) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = yaml.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, errors.New("YAML is evaluated as nil")
	}
	return v, nil
}

func Execute(inYaml io.Reader, out io.Writer, tmplFiles ...string) error {
	if len(tmplFiles) == 0 {
		return errors.New("no template files")
	}
	name := filepath.Base(tmplFiles[0])

	tmpl, err := template.New(name).Funcs(funcmap).ParseFiles(tmplFiles...)
	if err != nil {
		return err
	}

	v, err := readYAML(inYaml)
	if err != nil {
		return err
	}

	err = tmpl.Execute(out, v)
	if err != nil {
		return err
	}
	return nil
}
