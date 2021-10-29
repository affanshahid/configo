package configo

import (
	"encoding/json"

	"github.com/flynn/json5"
	"github.com/hjson/hjson-go"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type ProviderFunc func(in []byte) (map[string]interface{}, error)

func (p ProviderFunc) Parse(in []byte) (map[string]interface{}, error) {
	return p(in)
}

func parseYaml(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := yaml.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseJson(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := json.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseJson5(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := json5.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseHjson(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := hjson.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseToml(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := toml.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

var yamlProvider = ProviderFunc(parseYaml)
var jsonProvider = ProviderFunc(parseJson)
var json5Provider = ProviderFunc(parseJson5)
var hjsonProvider = ProviderFunc(parseHjson)
var tomlProvider = ProviderFunc(parseToml)
