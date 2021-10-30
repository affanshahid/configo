package configo_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/affanshahid/configo"
	"github.com/stretchr/testify/assert"
)

func ExampleConfig() {
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

	config, err := configo.NewConfig(dir, configo.WithDeploymentFromEnv("APP_ENV"))
	if err != nil {
		panic(err)
	}

	err = config.Initialize()
	if err != nil {
		panic(err)
	}

	fmt.Println(config.MustGetString("root.prop1"))
	fmt.Println(config.MustGetInt("root.prop2"))
	fmt.Println(config.MustGetBool("root.prop3"))
	fmt.Println(config.MustGetInt("root.prop4.nestedProp1"))

	// Output:
	// foo
	// 200
	// false
	// 4
}

func TestDefaultLoading(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                root:
                  prop1: foo
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "foo", config.MustGetString("root.prop1"))
}

func TestWithDeployment(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                root:
                    prop1: foo
                    prop2: bar
            `),
		},
		"production.yml": {
			Data: []byte(`
                root:
                    prop2: baz
            `),
		},
	}

	config, err := configo.NewConfig(dir, configo.WithDeployment("production"))
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "foo", config.MustGetString("root.prop1"))
	assert.Equal(t, "baz", config.MustGetString("root.prop2"))
}

func TestWithDeploymentFromEnv(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                root:
                    prop1: foo
                    prop2: bar
            `),
		},
		"production.yml": {
			Data: []byte(`
                root:
                    prop2: baz
            `),
		},
	}

	os.Setenv("DEP", "production")

	config, err := configo.NewConfig(dir, configo.WithDeploymentFromEnv("DEP"))
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "foo", config.MustGetString("root.prop1"))
	assert.Equal(t, "baz", config.MustGetString("root.prop2"))
}

func TestWithInstance(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                root:
                    prop1: foo
                    prop2: bar
            `),
		},
		"default-inst1.yml": {
			Data: []byte(`
                root:
                    prop2: baz
            `),
		},
	}

	config, err := configo.NewConfig(dir, configo.WithInstance("inst1"))
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "foo", config.MustGetString("root.prop1"))
	assert.Equal(t, "baz", config.MustGetString("root.prop2"))
}

func TestWithInstanceFromEnv(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                root:
                    prop1: foo
                    prop2: bar
            `),
		},
		"default-inst1.yml": {
			Data: []byte(`
                root:
                    prop2: baz
            `),
		},
	}

	os.Setenv("INST", "inst1")

	config, err := configo.NewConfig(dir, configo.WithInstanceFromEnv("INST"))
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "foo", config.MustGetString("root.prop1"))
	assert.Equal(t, "baz", config.MustGetString("root.prop2"))
}

func TestDefaultHostname(t *testing.T) {
	curFullHostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	curShortHostname := strings.Split(curFullHostname, ".")[0]

	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                root:
                    prop1: foo
                    prop2: bar
                    prop3: foo1
            `),
		},
		curFullHostname + ".yml": {
			Data: []byte(`
                root:
                    prop2: baz
            `),
		},
		curShortHostname + ".yml": {
			Data: []byte(`
                root:
                    prop3: foobar
            `),
		},
	}

	config, err := configo.NewConfig(dir, configo.WithInstance("inst1"))
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "foo", config.MustGetString("root.prop1"))
	assert.Equal(t, "baz", config.MustGetString("root.prop2"))
	assert.Equal(t, "foobar", config.MustGetString("root.prop3"))
}

func TestWithHostname(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                root:
                    prop1: foo
                    prop2: bar
                    prop3: foo1
            `),
		},
		"service1.yml": {
			Data: []byte(`
                root:
                    prop2: baz
            `),
		},
		"service1.example.com" + ".yml": {
			Data: []byte(`
                root:
                    prop3: foobar
            `),
		},
	}

	config, err := configo.NewConfig(dir, configo.WithHostname("service1.example.com"))
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "foo", config.MustGetString("root.prop1"))
	assert.Equal(t, "baz", config.MustGetString("root.prop2"))
	assert.Equal(t, "foobar", config.MustGetString("root.prop3"))
}

func TestWithHostnameFromEnv(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                root:
                    prop1: foo
                    prop2: bar
                    prop3: foo1
            `),
		},
		"service1.yml": {
			Data: []byte(`
                root:
                    prop2: baz
            `),
		},
		"service1.example.com" + ".yml": {
			Data: []byte(`
                root:
                    prop3: foobar
            `),
		},
	}
	os.Setenv("HOSTNAME", "service1.example.com")

	config, err := configo.NewConfig(dir, configo.WithHostnameFromEnv("HOSTNAME"))
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "foo", config.MustGetString("root.prop1"))
	assert.Equal(t, "baz", config.MustGetString("root.prop2"))
	assert.Equal(t, "foobar", config.MustGetString("root.prop3"))
}

func TestFileLoadingOrder(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p01: default
                p02: default
                p03: default
                p04: default
                p05: default
                p06: default
                p07: default
                p08: default
                p09: default
                p10: default
                p11: default
                p12: default
                p13: default
                p14: default
                p15: default
                p16: default
            `),
		},
		"default-inst1.yml": {
			Data: []byte(`
                p02: default-inst1
                p03: default-inst1
                p04: default-inst1
                p05: default-inst1
                p06: default-inst1
                p07: default-inst1
                p08: default-inst1
                p09: default-inst1
                p10: default-inst1
                p11: default-inst1
                p12: default-inst1
                p13: default-inst1
                p14: default-inst1
                p15: default-inst1
                p16: default-inst1
            `),
		},
		"production.yml": {
			Data: []byte(`
                p03: production
                p04: production
                p05: production
                p06: production
                p07: production
                p08: production
                p09: production
                p10: production
                p11: production
                p12: production
                p13: production
                p14: production
                p15: production
                p16: production
            `),
		},
		"production-inst1.yml": {
			Data: []byte(`
                p04: production-inst1
                p05: production-inst1
                p06: production-inst1
                p07: production-inst1
                p08: production-inst1
                p09: production-inst1
                p10: production-inst1
                p11: production-inst1
                p12: production-inst1
                p13: production-inst1
                p14: production-inst1
                p15: production-inst1
                p16: production-inst1
            `),
		},
		"service1.yml": {
			Data: []byte(`
                p05: service1
                p06: service1
                p07: service1
                p08: service1
                p09: service1
                p10: service1
                p11: service1
                p12: service1
                p13: service1
                p14: service1
                p15: service1
                p16: service1
            `),
		},
		"service1-inst1.yml": {
			Data: []byte(`
                p06: service1-inst1
                p07: service1-inst1
                p08: service1-inst1
                p09: service1-inst1
                p10: service1-inst1
                p11: service1-inst1
                p12: service1-inst1
                p13: service1-inst1
                p14: service1-inst1
                p15: service1-inst1
                p16: service1-inst1
            `),
		},
		"service1-production.yml": {
			Data: []byte(`
                p07: service1-production
                p08: service1-production
                p09: service1-production
                p10: service1-production
                p11: service1-production
                p12: service1-production
                p13: service1-production
                p14: service1-production
                p15: service1-production
                p16: service1-production
            `),
		},
		"service1-production-inst1.yml": {
			Data: []byte(`
                p08: service1-production-inst1
                p09: service1-production-inst1
                p10: service1-production-inst1
                p11: service1-production-inst1
                p12: service1-production-inst1
                p13: service1-production-inst1
                p14: service1-production-inst1
                p15: service1-production-inst1
                p16: service1-production-inst1
            `),
		},
		"service1.example.com.yml": {
			Data: []byte(`
                p09: service1.example.com
                p10: service1.example.com
                p11: service1.example.com
                p12: service1.example.com
                p13: service1.example.com
                p14: service1.example.com
                p15: service1.example.com
                p16: service1.example.com
            `),
		},
		"service1.example.com-inst1.yml": {
			Data: []byte(`
                p10: service1.example.com-inst1
                p11: service1.example.com-inst1
                p12: service1.example.com-inst1
                p13: service1.example.com-inst1
                p14: service1.example.com-inst1
                p15: service1.example.com-inst1
                p16: service1.example.com-inst1
            `),
		},
		"service1.example.com-production.yml": {
			Data: []byte(`
                p11: service1.example.com-production
                p12: service1.example.com-production
                p13: service1.example.com-production
                p14: service1.example.com-production
                p15: service1.example.com-production
                p16: service1.example.com-production
            `),
		},
		"service1.example.com-production-inst1.yml": {
			Data: []byte(`
                p12: service1.example.com-production-inst1
                p13: service1.example.com-production-inst1
                p14: service1.example.com-production-inst1
                p15: service1.example.com-production-inst1
                p16: service1.example.com-production-inst1
            `),
		},
		"local.yml": {
			Data: []byte(`
                p13: local
                p14: local
                p15: local
                p16: local
            `),
		},
		"local-inst1.yml": {
			Data: []byte(`
                p14: local-inst1
                p15: local-inst1
                p16: local-inst1
            `),
		},
		"local-production.yml": {
			Data: []byte(`
                p15: local-production
                p16: local-production
            `),
		},
		"local-production-inst1.yml": {
			Data: []byte(`
                p16: local-production-inst1
            `),
		},
	}

	config, err := configo.NewConfig(
		dir,
		configo.WithDeployment("production"),
		configo.WithInstance("inst1"),
		configo.WithHostname("service1.example.com"),
	)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "default", config.MustGetString("p01"))
	assert.Equal(t, "default-inst1", config.MustGetString("p02"))
	assert.Equal(t, "production", config.MustGetString("p03"))
	assert.Equal(t, "production-inst1", config.MustGetString("p04"))
	assert.Equal(t, "service1", config.MustGetString("p05"))
	assert.Equal(t, "service1-inst1", config.MustGetString("p06"))
	assert.Equal(t, "service1-production", config.MustGetString("p07"))
	assert.Equal(t, "service1-production-inst1", config.MustGetString("p08"))
	assert.Equal(t, "service1.example.com", config.MustGetString("p09"))
	assert.Equal(t, "service1.example.com-inst1", config.MustGetString("p10"))
	assert.Equal(t, "service1.example.com-production", config.MustGetString("p11"))
	assert.Equal(t, "service1.example.com-production-inst1", config.MustGetString("p12"))
	assert.Equal(t, "local", config.MustGetString("p13"))
	assert.Equal(t, "local-inst1", config.MustGetString("p14"))
	assert.Equal(t, "local-production", config.MustGetString("p15"))
	assert.Equal(t, "local-production-inst1", config.MustGetString("p16"))
}

func TestGet(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: foo
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.Get("p1")
	assert.Nil(t, err)
	assert.Equal(t, "foo", val)
}

func TestGetMissingPropertyError(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: foo
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	_, err = config.Get("p2.p3")
	assert.NotNil(t, err)
}

func TestGetString(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: foo
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetString("p1")
	assert.Nil(t, err)
	assert.Equal(t, "foo", val)
}

func TestGetBool(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: true
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetBool("p1")
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestGetInt(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: 10
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetInt("p1")
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func TestGetInt32(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: 10
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetInt32("p1")
	assert.Nil(t, err)
	assert.Equal(t, int32(10), val)
}

func TestGetInt64(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: 10
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetInt64("p1")
	assert.Nil(t, err)
	assert.Equal(t, int64(10), val)
}

func TestGetUint(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: 10
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetUint("p1")
	assert.Nil(t, err)
	assert.Equal(t, uint(10), val)
}

func TestGetUint32(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: 10
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetUint32("p1")
	assert.Nil(t, err)
	assert.Equal(t, uint32(10), val)
}

func TestGetUint64(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: 10
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetUint64("p1")
	assert.Nil(t, err)
	assert.Equal(t, uint64(10), val)
}

func TestGetFloat64(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: 10.0
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetFloat64("p1")
	assert.Nil(t, err)
	assert.Equal(t, 10.0, val)
}

func TestGetTime(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: 1635565664
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetTime("p1")
	assert.Nil(t, err)
	assert.Equal(t, time.Unix(1635565664, 0), val)
}

func TestGetDuration(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1: 10h
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetDuration("p1")
	assert.Nil(t, err)
	assert.Equal(t, time.Duration(10*time.Hour), val)
}

func TestGetIntSlice(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1:
                  - 1
                  - 2
                  - 3
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetIntSlice("p1")
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 3}, val)
}

func TestGetStringSlice(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1:
                  - foo
                  - bar
                  - baz
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetStringSlice("p1")
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo", "bar", "baz"}, val)
}

func TestGetStringMap(t *testing.T) {
	dir := fstest.MapFS{
		"default.yml": {
			Data: []byte(`
                p1:
                  foo: 1
                  bar: 2
                  baz: 3
            `),
		},
	}

	config, err := configo.NewConfig(dir)
	assert.Nil(t, err, "err should be nil")

	err = config.Initialize()
	assert.Nil(t, err, "err should be nil")

	val, err := config.GetStringMap("p1")
	assert.Nil(t, err)
	assert.Equal(
		t,
		map[string]interface{}{
			"foo": 1,
			"bar": 2,
			"baz": 3,
		},
		val,
	)
}
