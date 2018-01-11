package tmpl

import (
	"bytes"

	yaml "gopkg.in/yaml.v2"
)

func init() {
	AddFunc("yaml", funcYAML)
}

var lf = []byte{'\n'}

func indentYAML(b []byte, indent string) []byte {
	for b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	x := append(lf, []byte(indent)...)
	return bytes.Replace(b, lf, x, -1)
}

func funcYAML(v interface{}, indent string) (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(indentYAML(b, indent)), nil
}
