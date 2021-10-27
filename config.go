package configo

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/imdario/mergo"
)

type environment struct {
	deployment string
}

type Config struct {
	environment
	dir   string
	store map[string]interface{}
}

type ConfigOption func(*Config)

func NewConfig(dir string, opts ...ConfigOption) *Config {
	c := &Config{dir: dir, environment: environment{"development"}}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithDeployment(deployment string) ConfigOption {
	return func(c *Config) {
		c.deployment = deployment
	}
}

func (c *Config) Initialize() error {
	c.store = map[string]interface{}{}

	files, err := os.ReadDir(c.dir)
	if err != nil {
		return err
	}

	fileMap := map[string]fs.DirEntry{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		nameWithoutExt := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		fileMap[nameWithoutExt] = file
	}

	for _, tmpl := range orderedTemplates {
		filename := getExpectedBasename(tmpl, c.environment)

		entry, found := fileMap[filename]
		if !found {
			continue
		}

		path := filepath.Join(c.dir, entry.Name())
		in, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		provider := defaultProviders[ext]
		data, err := provider.Parse(in)
		if err != nil {
			return err
		}

		err = mergo.Merge(&c.store, data, mergo.WithOverride)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) TryGet(path string) (interface{}, error) {
	return jsonpath.Get(path, c.store)
}

func (c *Config) Get(path string) interface{} {
	v, err := jsonpath.Get(path, c.store)
	if err != nil {
		panic(err)
	}
	return v
}
