package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ConfigureRoutes(app *fiber.App, db *gorm.DB) {
	app.Static("/", "./public")
	api := app.Group("/api")
	api.Get("/idea", ListIdeas(db))
	api.Post("/idea", CreateIdea(db))
	api.Get("/idea/:id", GetIdea(db))
	api.Put("/idea/:id", UpdateIdea(db))
	api.Delete("/idea/:id", DeleteIdea(db))
	api.Post("/idea/:id/comment", CreateIdeaComment(db))
	api.Put("/idea/:id/comment", UpdateIdeaComment(db))
	api.Delete("/idea/:id/comment", DeleteIdeaComment(db))
}
