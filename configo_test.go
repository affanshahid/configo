package configo_test

import (
	"fmt"
	"os"
	"testing"
	"testing/fstest"

	"github.com/affanshahid/configo"
	"github.com/stretchr/testify/assert"
)

func Example() {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                root:
                  prop1: foo
                  prop2: 100
                  prop3: false
                  prop4:
                    nestedProp1: 4
            `),
		},
		"production.yml": {
			Data: []byte(`
                root:
                  prop2: 200
            `),
		},
	}

	os.Setenv("APP_ENV", "production")

	err := configo.Initialize(dir, configo.WithDeploymentFromEnv("APP_ENV"))
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

func TestGlobalConfig(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                root:
                  prop1: foo
            `),
		},
	}

	err := configo.Initialize(dir)
	assert.Nilf(t, err, "err should be nil")

	assert.Equal(t, "foo", configo.MustGetString("root.prop1"))
}

func TestGlobalConfigParseError(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                in valid
                ; hhvc: foo
            `),
		},
	}

	err := configo.Initialize(dir)
	assert.NotNil(t, err)
}
