package configo

import (
	"io/fs"
)

type Config struct {
	dir fs.ReadDirFile

	deployment string
}

type ConfigOption func(*Config)

func NewConfig(dir fs.ReadDirFile, opts ...ConfigOption) *Config {
	c := &Config{dir: dir}

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
