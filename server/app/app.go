package app

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ideashare/config"
	"ideashare/models"
	"ideashare/routes"
)

func RunApp(dbDsn string) *fiber.App {
	db, err := gorm.Open(mysql.Open(dbDsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err = db.AutoMigrate(
		&models.User{},
		&models.Idea{},
		&models.IdeaComment{},
		&models.IdeaLike{},
		&models.SiteSetting{},
		&models.UserRole{},
		&models.UserSetting{},
	); err != nil {
		panic(err)
	}
	app := fiber.New()
	app.Server().WriteBufferSize = 1024 * 1024 * 1024
	app.Server().ReadBufferSize = 256 * 1024
	routes.ConfigureRoutes(app, &config.AppContainer{Db: db})
	print("Starting server...")

	return app
}
