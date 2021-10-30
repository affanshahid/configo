# Configo

A Go port of [node-config](https://github.com/lorenwest/node-config).

[Documentation](https://pkg.go.dev/github.com/affanshahid/configo)

## Installing

```sh
go get -u github.com/affanshahid/configo
```

## Usage

```yml
# config/default.yml
root:
  prop1: foo
  prop2: 100
  prop3: false
  prop4:
    nestedProp1: 4
```

```yml
# config/production.yml
root:
  prop2: 200
```

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
