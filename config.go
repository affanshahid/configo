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
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/affanshahid/walkmap"
	"github.com/imdario/mergo"
	"github.com/spf13/cast"
)

const development = "development"

type environment struct {
	deployment, instance, shortHostname, fullHostname string
}

type Config struct {
	environment
	dir   string
	store map[string]interface{}
}

type ConfigOption func(*Config)

func NewConfig(dir string, opts ...ConfigOption) (*Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	shortHostname := strings.Split(hostname, ".")[0]

	c := &Config{dir: dir, environment: environment{development, "", shortHostname, hostname}}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func WithDeployment(deployment string) ConfigOption {
	return func(c *Config) {
		c.deployment = deployment
	}
}

func WithInstance(instance string) ConfigOption {
	return func(c *Config) {
		c.instance = instance
	}
}

func WithHostname(hostname string) ConfigOption {
	return func(c *Config) {
		c.shortHostname = strings.Split(hostname, ".")[0]
		c.fullHostname = hostname
	}
}

func WithDeploymentFromEnv(env string) ConfigOption {
	deployment, exists := os.LookupEnv(env)
	if exists {
		return WithDeployment(deployment)
	} else {
		return WithDeployment(development)
	}
}

func WithInstanceFromEnv(env string) ConfigOption {
	return WithInstance(os.Getenv(env))
}

func WithHostnameFromEnv(env string) ConfigOption {
	return WithHostname(os.Getenv(env))
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
		if filename == "" {
			continue
		}

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
	for i := 0; i < len(paths)-1; i++ {
		s := paths[i]
		if v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		if i, err := strconv.Atoi(s); err == nil {
			v = v.Index(i)
		}
		v = v.MapIndex(reflect.ValueOf(s))
	}

	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	v.SetMapIndex(reflect.ValueOf(paths[len(paths)-1]), reflect.ValueOf(val))
}

func (c *Config) Get(path string) (interface{}, error) {
	return jsonpath.Get(path, c.store)
}

func (c *Config) GetString(path string) (string, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return "", err
	}

	return cast.ToStringE(out)
}

func (c *Config) GetBool(path string) (bool, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return false, err
	}

	return cast.ToBoolE(out)
}

func (c *Config) GetInt(path string) (int, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToIntE(out)
}

func (c *Config) GetInt32(path string) (int32, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToInt32E(out)
}

func (c *Config) GetInt64(path string) (int64, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToInt64E(out)
}

func (c *Config) GetUint(path string) (uint, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToUintE(out)
}

func (c *Config) GetUint32(path string) (uint32, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToUint32E(out)
}

func (c *Config) GetUint64(path string) (uint64, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToUint64E(out)
}

func (c *Config) GetFloat64(path string) (float64, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToFloat64E(out)
}

func (c *Config) GetTime(path string) (time.Time, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return time.Time{}, err
	}

	return cast.ToTimeE(out)
}

func (c *Config) GetDuration(path string) (time.Duration, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return time.Duration(0), err
	}

	return cast.ToDurationE(out)
}

func (c *Config) GetIntSlice(path string) ([]int, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return nil, err
	}

	return cast.ToIntSliceE(out)
}

func (c *Config) GetStringSlice(path string) ([]string, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return nil, err
	}

	return cast.ToStringSliceE(out)
}

func (c *Config) GetStringMap(path string) (map[string]interface{}, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return nil, err
	}

	return cast.ToStringMapE(out)
}

func (c *Config) MustGet(path string) interface{} {
	v, err := c.Get(path)
	if err != nil {
		panic(err)
	}
	return v
}
