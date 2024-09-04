package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateIdeaComment(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func UpdateIdeaComment(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func DeleteIdeaComment(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
