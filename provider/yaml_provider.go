package provider

import (
	"gopkg.in/yaml.v2"
)

type YamlProvider struct{}

func (p *YamlProvider) Parse(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := yaml.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
