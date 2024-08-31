package config

import (
	"fmt"
	"os"
)

var overrides = make(map[string]string)

const (
	DbHost = "DB_HOST"
	DbPort = "DB_PORT"
	DbUser = "DB_USER"
	DbPass = "DB_PASS"
	DbName = "DB_NAME"
)

func GetStringOr(key string, def string) string {
	if val, ok := overrides[key]; ok {
		return val
	}
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return val
}

func MustGetString(key string) string {
	if val, ok := overrides[key]; ok {
		return val
	}
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return val
}

func SetOverride(key string, value string) {
	overrides[key] = value
}
