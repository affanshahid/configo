package configo

import (
	"io/fs"
	"time"
)

var globalConfig *Config

// Initialize loads the global configuration
func Initialize(dir fs.FS, opts ...ConfigOption) (err error) {
	globalConfig, err = NewConfig(dir, opts...)
	if err != nil {
		return err
	}

	err = globalConfig.Initialize()
	if err != nil {
		return err
	}

	return nil
}

// Get returns the value at the given path as an interface from the globalConfig
func Get(path string) (interface{}, error) {
	return globalConfig.Get(path)
}

// GetString returns the value at the given path as a string from the globalConfig
func GetString(path string) (string, error) {
	return globalConfig.GetString(path)
}

// GetBool returns the value at the given path as a boolean from the globalConfig
func GetBool(path string) (bool, error) {
	return globalConfig.GetBool(path)
}

// GetInt returns the value at the given path as a int from the globalConfig
func GetInt(path string) (int, error) {
	return globalConfig.GetInt(path)
}

// GetInt32 returns the value at the given path as a int32 from the globalConfig
func GetInt32(path string) (int32, error) {
	return globalConfig.GetInt32(path)
}

// GetInt64 returns the value at the given path as a int64 from the globalConfig
func GetInt64(path string) (int64, error) {
	return globalConfig.GetInt64(path)
}

// GetUint returns the value at the given path as a uint from the globalConfig
func GetUint(path string) (uint, error) {
	return globalConfig.GetUint(path)
}

// GetUint32 returns the value at the given path as a uint32 from the globalConfig
func GetUint32(path string) (uint32, error) {
	return globalConfig.GetUint32(path)
}

// GetUint64 returns the value at the given path as a uint64 from the globalConfig
func GetUint64(path string) (uint64, error) {
	return globalConfig.GetUint64(path)
}

// GetFloat64 returns the value at the given path as a float64 from the globalConfig
func GetFloat64(path string) (float64, error) {
	return globalConfig.GetFloat64(path)
}

// GetTime returns the value at the given path as time
func GetTime(path string) (time.Time, error) {
	return globalConfig.GetTime(path)
}

// GetDuration returns the value at the given path as a duration from the globalConfig
func GetDuration(path string) (time.Duration, error) {
	return globalConfig.GetDuration(path)
}

// GetIntSlice returns the value at the given path as a slice of int values from the globalConfig
func GetIntSlice(path string) ([]int, error) {
	return globalConfig.GetIntSlice(path)
}

// GetStringSlice returns the value at the given path as a slice of string values from the globalConfig
func GetStringSlice(path string) ([]string, error) {
	return globalConfig.GetStringSlice(path)
}

// GetStringMap returns the value at the given path as a map with string keys from the globalConfig
// and values as interfaces
func GetStringMap(path string) (map[string]interface{}, error) {
	return globalConfig.GetStringMap(path)
}

// MustGet is the same as `Get` except it panics in case of an error
func MustGet(path string) interface{} {
	return globalConfig.MustGet(path)
}

// MustGetString is the same as `GetString` except it panics in case of an error
func MustGetString(path string) string {
	return globalConfig.MustGetString(path)
}

// MustGetBool is the same as `GetBool` except it panics in case of an error
func MustGetBool(path string) bool {
	return globalConfig.MustGetBool(path)
}

// MustGetInt is the same as `GetInt` except it panics in case of an error
func MustGetInt(path string) int {
	return globalConfig.MustGetInt(path)
}

// MustGetInt32 is the same as `GetInt32` except it panics in case of an error
func MustGetInt32(path string) int32 {
	return globalConfig.MustGetInt32(path)
}

// MustGetInt64 is the same as `GetInt64` except it panics in case of an error
func MustGetInt64(path string) int64 {
	return globalConfig.MustGetInt64(path)
}

// MustGetUint is the same as `GetUint` except it panics in case of an error
func MustGetUint(path string) uint {
	return globalConfig.MustGetUint(path)
}

// MustGetUint32 is the same as `GetUint32` except it panics in case of an error
func MustGetUint32(path string) uint32 {
	return globalConfig.MustGetUint32(path)
}

// MustGetUint64 is the same as `GetUint64` except it panics in case of an error
func MustGetUint64(path string) uint64 {
	return globalConfig.MustGetUint64(path)
}

// MustGetFloat64 is the same as `GetFloat64` except it panics in case of an error
func MustGetFloat64(path string) float64 {
	return globalConfig.MustGetFloat64(path)
}

// MustGetTime is the same as `GetTime` except it panics in case of an error
func MustGetTime(path string) time.Time {
	return globalConfig.MustGetTime(path)
}

// MustGetDuration is the same as `GetDuration` except it panics in case of an error
func MustGetDuration(path string) time.Duration {
	return globalConfig.MustGetDuration(path)
}

// MustGetIntSlice is the same as `GetIntSlice` except it panics in case of an error
func MustGetIntSlice(path string) []int {
	return globalConfig.MustGetIntSlice(path)
}

// MustGetStringSlice is the same as `GetStringSlice` except it panics in case of an error
func MustGetStringSlice(path string) []string {
	return globalConfig.MustGetStringSlice(path)
}

// MustGetStringMap is the same as `GetStringMap` except it panics in case of an error
func MustGetStringMap(path string) map[string]interface{} {
	return globalConfig.MustGetStringMap(path)
}
