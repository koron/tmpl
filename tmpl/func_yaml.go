package tmpl

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func init() {
	AddFunc("yaml", funcYAML)
}

func funcYAML(v interface{}, indent string) (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
