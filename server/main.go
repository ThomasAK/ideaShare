package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ideashare/config"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.GetStringOr(config.DbUser, "ideashare"),
		config.GetStringOr(config.DbPass, "password"),
		config.GetStringOr(config.DbHost, "localhost"),
		config.GetStringOr(config.DbPort, "3318"),
		config.GetStringOr(config.DbName, "ideashare"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := goose.SetDialect("mysql"); err != nil {
		panic(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	if err := goose.Up(sqlDb, "migrations"); err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Static("/", "./public")
	print("Starting server...")

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
