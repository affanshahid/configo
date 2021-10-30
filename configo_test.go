package configo_test

import (
	"testing"
	"testing/fstest"

	"github.com/affanshahid/configo"
	"github.com/stretchr/testify/assert"
)

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
