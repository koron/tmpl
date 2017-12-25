package tmpl

import yaml "gopkg.in/yaml.v2"

func init() {
	AddFunc("yaml", funcYAML)
}

func funcYAML(v interface{}, indent string) ([]byte, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}
