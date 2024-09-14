package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ideashare/config"
	"ideashare/models"
	"ideashare/routes"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
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
	if err = db.AutoMigrate(
		&models.Idea{},
		&models.IdeaComment{},
		&models.IdeaLike{},
		&models.SiteSetting{},
		&models.User{},
		&models.UserRole{},
		&models.UserSetting{},
	); err != nil {
		panic(err)
	}
	app := fiber.New()
	app.Server().ReadBufferSize = 256 * 1024
	routes.ConfigureRoutes(app, &config.AppContainer{Db: db})
	print("Starting server...")

	if err := app.Listen(":3030"); err != nil {
		panic(err)
	}
}
