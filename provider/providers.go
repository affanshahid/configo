package provider

import (
	"encoding/json"

	"github.com/flynn/json5"
	"github.com/hjson/hjson-go"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
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

type JsonProvider struct{}

func (p *JsonProvider) Parse(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := json.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type Json5Provider struct{}

func (p *Json5Provider) Parse(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := json5.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type HjsonProvider struct{}

func (p *HjsonProvider) Parse(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := hjson.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type TomlProvider struct{}

func (p *TomlProvider) Parse(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := toml.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
