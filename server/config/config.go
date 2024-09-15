package config

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
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

func MakeDsn(user, pass, host, port, name string, allowNativePasswords bool) string {
	c := mysql.Config{
		Addr:                 host + ":" + port,
		Net:                  "tcp",
		DBName:               name,
		ParseTime:            true,
		Loc:                  time.UTC,
		AllowNativePasswords: allowNativePasswords,
	}
	if user != "" {
		c.User = user
		c.Passwd = pass
	}
	return c.FormatDSN()
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
