package config

import (
	"fmt"
	"gorm.io/gorm"
	"os"
)

var overrides = make(map[string]string)

type AppContainer struct {
	Db *gorm.DB
}

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

func MakeDsn(user, pass, host, port, name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user,
		pass,
		host,
		port,
		name,
	)
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
