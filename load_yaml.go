package tmpl

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// SourceYAML is a source for LoadYAML function. Default is os.Stdin
var SourceYAML io.Reader

// LoadYAML loads a YAML from SourceYAML.
func LoadYAML() (interface{}, error) {
	r := SourceYAML
	if SourceYAML == nil {
		r = os.Stdin
	}
	return readYAML(r)
}

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
