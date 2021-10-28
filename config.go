package configo

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/affanshahid/walkmap"
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

		data, err := c.readFile(entry.Name())
		if err != nil {
			return err
		}

		err = mergo.Merge(&c.store, data, mergo.WithOverride)
		if err != nil {
			return err
		}
	}

	if envFile, found := fileMap[envFileName]; found {
		data, err := c.readFile(envFile.Name())
		if err != nil {
			return err
		}

		err = c.loadOverrides(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) readFile(name string) (map[string]interface{}, error) {
	path := filepath.Join(c.dir, name)
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(name))
	provider := defaultProviders[ext]
	return provider.Parse(in)
}

func (c *Config) loadOverrides(data map[string]interface{}) error {
	var err error
	walkmap.Walk(data, func(keyPath []interface{}, value interface{}, kind reflect.Kind) {
		if err != nil {
			return
		}

		strPath := make([]string, len(keyPath))

		for i, p := range keyPath {
			strPath[i] = p.(string)
		}

		envName, ok := value.(string)
		if !ok {
			err = fmt.Errorf("invalid environment variable %s", envName)
			return
		}

		if envValue, found := os.LookupEnv(envName); found {
			c.set(strPath, envValue)
		}
	})

	return err
}

func (c *Config) set(paths []string, val interface{}) {
	v := reflect.ValueOf(c.store)
	for _, s := range paths {
		l := reflect.TypeOf(v.Interface())
		k := l.Kind()

		l2 := reflect.TypeOf(s)
		k2 := l2.Kind()
		t := v.Interface()
		fmt.Println(k, k2, t)

		if i, err := strconv.Atoi(s); err == nil {
			v = v.Index(i)
		}
		v = v.MapIndex(reflect.ValueOf(s))
	}
	v.Set(reflect.ValueOf(val))
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
