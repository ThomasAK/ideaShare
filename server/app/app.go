package app

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ideashare/config"
	"ideashare/models"
	"ideashare/routes"
	"time"
)

func RunApp(dbDsn string) (*fiber.App, *config.AppContainer) {
	time.Local = time.UTC
	db, err := gorm.Open(mysql.Open(dbDsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err = db.AutoMigrate(
		&models.User{},
		&models.UserRole{},
		&models.UserSetting{},
		&models.Idea{},
		&models.IdeaLike{},
		&models.IdeaComment{},
		&models.SiteSetting{},
	); err != nil {
		panic(err)
	}
	app := fiber.New()
	app.Server().WriteBufferSize = 1024 * 1024 * 1024
	app.Server().ReadBufferSize = 256 * 1024
	container := &config.AppContainer{Db: db}
	routes.ConfigureRoutes(app, container)
	print("Starting server...")

	return app, container
}
