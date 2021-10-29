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

const development = "dev"

type environment struct {
	deployment, instance, shortHostname, fullHostname string
}

// Config is a hierarchical loader and access point for configurations.
// It allows loading configurations from mutiple files while being cognizant
// of the einvironment.
// Files are loading in the following order:
//	 default.EXT
//	 default-{instance}.EXT
//	 {deployment}.EXT
//	 {deployment}-{instance}.EXT
//	 {short_hostname}.EXT
//	 {short_hostname}-{instance}.EXT
//	 {short_hostname}-{deployment}.EXT
//	 {short_hostname}-{deployment}-{instance}.EXT
//	 {full_hostname}.EXT
//	 {full_hostname}-{instance}.EXT
//	 {full_hostname}-{deployment}.EXT
//	 {full_hostname}-{deployment}-{instance}.EXT
//	 local.EXT
//	 local-{instance}.EXT
//	 local-{deployment}.EXT
//	 local-{deployment}-{instance}.EXT
//
// EXT can be: `yaml`, `yml`, `json`, `json5`, `hjson`, `toml`
//
// deployment defines your current environment i.e dev, test, prod etc (defaults to "dev")
//
// instance can be the node ID in a multi-node deployment (defaults to "")
//
// shortHostname is the hostname till the first `.` (derived from `os.Hostname()` by default)
//
// fullHostname is the full host name (defaults to `os.Hostname()`)
//
// Each file overrides configurations from the file above.
// There is a special file called `env.EXT` which allows overriding
// configurations using environment variables
type Config struct {
	environment
	dir   string
	store map[string]interface{}
}

// ConfigOption is a functional option to configure a Config instance
type ConfigOption func(*Config)

// NewConfig creates a new Config
// dir is the path to the directory containing the config files
// opts are functional options
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

// WithDeployment sets the given deployment
func WithDeployment(deployment string) ConfigOption {
	return func(c *Config) {
		c.deployment = deployment
	}
}

// WithInstance sets the given instance
func WithInstance(instance string) ConfigOption {
	return func(c *Config) {
		c.instance = instance
	}
}

// WithHostname uses the given string to set shortHostname and fullHostname
func WithHostname(hostname string) ConfigOption {
	return func(c *Config) {
		c.shortHostname = strings.Split(hostname, ".")[0]
		c.fullHostname = hostname
	}
}

// WithDeploymentFromEnv loads the deployment label from the given environment variable
func WithDeploymentFromEnv(env string) ConfigOption {
	deployment, exists := os.LookupEnv(env)
	if exists {
		return WithDeployment(deployment)
	} else {
		return WithDeployment(development)
	}
}

// WithInstanceFromEnv loads the instance id from the given environment variable
func WithInstanceFromEnv(env string) ConfigOption {
	return WithInstance(os.Getenv(env))
}

// WithHostnameFromEnv loads the hostname from the given environment variable
func WithHostnameFromEnv(env string) ConfigOption {
	return WithHostname(os.Getenv(env))
}

// Initialize initializes and loads in the configurations
// This must be called before attempting to get values
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

// Get returns the value at the given path as an interface
func (c *Config) Get(path string) (interface{}, error) {
	return jsonpath.Get(path, c.store)
}

// GetString returns the value at the given path as a string
func (c *Config) GetString(path string) (string, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return "", err
	}

	return cast.ToStringE(out)
}

// GetBool returns the value at the given path as a boolean
func (c *Config) GetBool(path string) (bool, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return false, err
	}

	return cast.ToBoolE(out)
}

// GetInt returns the value at the given path as a int
func (c *Config) GetInt(path string) (int, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToIntE(out)
}

// GetInt32 returns the value at the given path as a int32
func (c *Config) GetInt32(path string) (int32, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToInt32E(out)
}

// GetInt64 returns the value at the given path as a int64
func (c *Config) GetInt64(path string) (int64, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToInt64E(out)
}

// GetUint returns the value at the given path as a uint
func (c *Config) GetUint(path string) (uint, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToUintE(out)
}

// GetUint32 returns the value at the given path as a uint32
func (c *Config) GetUint32(path string) (uint32, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToUint32E(out)
}

// GetUint64 returns the value at the given path as a uint64
func (c *Config) GetUint64(path string) (uint64, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToUint64E(out)
}

// GetFloat64 returns the value at the given path as a float64
func (c *Config) GetFloat64(path string) (float64, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return 0, err
	}

	return cast.ToFloat64E(out)
}

// GetTime returns the value at the given path as time
func (c *Config) GetTime(path string) (time.Time, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return time.Time{}, err
	}

	return cast.ToTimeE(out)
}

// GetDuration returns the value at the given path as a duration
func (c *Config) GetDuration(path string) (time.Duration, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return time.Duration(0), err
	}

	return cast.ToDurationE(out)
}

// GetIntSlice returns the value at the given path as a slice of int values
func (c *Config) GetIntSlice(path string) ([]int, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return nil, err
	}

	return cast.ToIntSliceE(out)
}

// GetStringSlice returns the value at the given path as a slice of string values
func (c *Config) GetStringSlice(path string) ([]string, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return nil, err
	}

	return cast.ToStringSliceE(out)
}

// GetStringMap returns the value at the given path as a map with string keys
// and values as interfaces
func (c *Config) GetStringMap(path string) (map[string]interface{}, error) {
	out, err := jsonpath.Get(path, c.store)
	if err != nil {
		return nil, err
	}

	return cast.ToStringMapE(out)
}

// MustGet is the same as `Get` except it panics in case of an error
func (c *Config) MustGet(path string) interface{} {
	v, err := c.Get(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetString is the same as `GetString` except it panics in case of an error
func (c *Config) MustGetString(path string) string {
	v, err := c.GetString(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetBool is the same as `GetBool` except it panics in case of an error
func (c *Config) MustGetBool(path string) bool {
	v, err := c.GetBool(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetInt is the same as `GetInt` except it panics in case of an error
func (c *Config) MustGetInt(path string) int {
	v, err := c.GetInt(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetInt32 is the same as `GetInt32` except it panics in case of an error
func (c *Config) MustGetInt32(path string) int32 {
	v, err := c.GetInt32(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetInt64 is the same as `GetInt64` except it panics in case of an error
func (c *Config) MustGetInt64(path string) int64 {
	v, err := c.GetInt64(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetUint is the same as `GetUint` except it panics in case of an error
func (c *Config) MustGetUint(path string) uint {
	v, err := c.GetUint(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetUint32 is the same as `GetUint32` except it panics in case of an error
func (c *Config) MustGetUint32(path string) uint32 {
	v, err := c.GetUint32(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetUint64 is the same as `GetUint64` except it panics in case of an error
func (c *Config) MustGetUint64(path string) uint64 {
	v, err := c.GetUint64(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetFloat64 is the same as `GetFloat64` except it panics in case of an error
func (c *Config) MustGetFloat64(path string) float64 {
	v, err := c.GetFloat64(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetTime is the same as `GetTime` except it panics in case of an error
func (c *Config) MustGetTime(path string) time.Time {
	v, err := c.GetTime(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetDuration is the same as `GetDuration` except it panics in case of an error
func (c *Config) MustGetDuration(path string) time.Duration {
	v, err := c.GetDuration(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetIntSlice is the same as `GetIntSlice` except it panics in case of an error
func (c *Config) MustGetIntSlice(path string) []int {
	v, err := c.GetIntSlice(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetStringSlice is the same as `GetStringSlice` except it panics in case of an error
func (c *Config) MustGetStringSlice(path string) []string {
	v, err := c.GetStringSlice(path)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetStringMap is the same as `GetStringMap` except it panics in case of an error
func (c *Config) MustGetStringMap(path string) map[string]interface{} {
	v, err := c.GetStringMap(path)
	if err != nil {
		panic(err)
	}
	return v
}
