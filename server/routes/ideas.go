package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"ideashare/models"
)

func ListIdeas(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var ideas []models.Idea
		result := db.Find(&ideas)
		if result.Error != nil {
			return result.Error
		}
		if err := ctx.JSON(ideas); err != nil {
			return ctx.SendStatus(500)
		}
		return nil
	}
}

func CreateIdea(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		idea := &models.Idea{}
		if err := ctx.BodyParser(idea); err != nil {
			return ctx.SendStatus(400)
		}
		idea.CreatedBy = 1
		result := db.Create(idea)
		if result.Error != nil {
			return result.Error
		}
		ctx.Status(201)
		if err := ctx.JSON(idea); err != nil {
			return ctx.SendStatus(500)
		}
		return nil
	}
}

func GetIdea(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func UpdateIdea(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func DeleteIdea(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
