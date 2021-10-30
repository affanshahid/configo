# Configo

A Go port of [node-config](https://github.com/lorenwest/node-config).

[Documentation](https://pkg.go.dev/github.com/affanshahid/configo)

## Introduction

Configo enables hierarchical configurations for your application.

It allows you to setup default configuration parameters and override them based on
the environment the application is running in. (i.e development, staging, production etc)

Configurations can be stored in a number of formats (yaml, json, json5, hjson and toml). They can also be overriden using custom environment variables.

There are a number of parameters which decide which configuration files are loaded at runtime.

Files are loading in the following order (Each file is optional, you may only create ones you expect to use):

- `default.EXT`
- `default-{instance}.EXT`
- `{deployment}.EXT`
- `{deployment}-{instance}.EXT`
- `{short_hostname}.EXT`
- `{short_hostname}-{instance}.EXT`
- `{short_hostname}-{deployment}.EXT`
- `{short_hostname}-{deployment}-{instance}.EXT`
- `{full_hostname}.EXT`
- `{full_hostname}-{instance}.EXT`
- `{full_hostname}-{deployment}.EXT`
- `{full_hostname}-{deployment}-{instance}.EXT`
- `local.EXT`
- `local-{instance}.EXT`
- `local-{deployment}.EXT`
- `local-{deployment}-{instance}.EXT`

EXT can be: `yaml`, `yml`, `json`, `json5`, `hjson`, `toml`

`deployment` defines your current environment i.e dev, test, prod etc (defaults to "dev")

`instance` can be the node ID in a multi-node deployment (defaults to "")

`shortHostname` is the hostname till the first `.` (derived from `os.Hostname()` by default)

`fullHostname` is the full host name (defaults to `os.Hostname()`)

Each file overrides configurations from the file above.
There is a special file called `env.EXT` which allows overriding
configurations using environment variables

## Installing

```sh
go get -u github.com/affanshahid/configo
```

## Usage

First create a configuration file to define your default parameters

```yml
# config/default.yml
root:
  prop1: foo
  prop2: 100
  prop3: false
  prop4:
    nestedProp1: 4
```

Then create a file to define the profile you want to use when in production. Here we only override one parameter and expect the defaults for the others

```yml
# config/production.yml
root:
  prop2: 200
```

You may also create a file to optionally override all parameters from environment variables

```yml
# config/env.yml
root:
  prop1: PROP1_ENV
  prop2: PROP2_ENV
  prop3: PROP3_ENV
  prop4:
    nestedProp1: NESTED_PROP1_ENV
```

Set an environment variable to specify the current deployment (dev, production, staging etc)

```sh
export APP_ENV=production
```

```go
package main

import (
	"fmt"
	"os"

	"github.com/affanshahid/configo"
)

func main() {
	// Here we initialize the configurations and specify
	// the environment variable from which to load the
	// deployment name (the one we set above)
	err := configo.Initialize(os.DirFS("./config"), configo.WithDeploymentFromEnv("APP_ENV"))
	if err != nil {
		panic(err)
	}

	fmt.Println(configo.MustGetString("root.prop1"))
	fmt.Println(configo.MustGetInt("root.prop2"))
	fmt.Println(configo.MustGetBool("root.prop3"))
	fmt.Println(configo.MustGetInt("root.prop4.nestedProp1"))

	// Output:
	// foo
	// 200
	// false
	// 4
}
```

## See Also

- [node-config](https://github.com/lorenwest/node-config) (Main inspiration)
