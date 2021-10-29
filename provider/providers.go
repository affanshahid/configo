package provider

import (
	"encoding/json"

	"github.com/flynn/json5"
	"github.com/hjson/hjson-go"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// YamlProvider is a `Provider` for yml files
type YamlProvider struct{}

// Parse parses yml content
func (p *YamlProvider) Parse(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := yaml.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// JsonProvider is a `Provider` for json files
type JsonProvider struct{}

// Parse parses json content
func (p *JsonProvider) Parse(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := json.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Json5Provider is a `Provider` for json5 files
type Json5Provider struct{}

// Parse parses json5 content
func (p *Json5Provider) Parse(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := json5.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// HjsonProvider is a `Provider` for hjson files
type HjsonProvider struct{}

// Parse parses hjson content
func (p *HjsonProvider) Parse(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := hjson.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// TomlProvider is a `Provider` for toml files
type TomlProvider struct{}

// Parse parses toml content
func (p *TomlProvider) Parse(in []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := toml.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
