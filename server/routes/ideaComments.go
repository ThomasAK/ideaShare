package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"ideashare/models"
)

func CreateIdeaComment(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ideaComment := &models.IdeaComment{}
		if err := ctx.BodyParser(ideaComment); err != nil {
			return ctx.SendStatus(400)
		}
		ideaComment.CreatedBy = 1
		result := db.Create(ideaComment)
		if result.Error != nil {
			return result.Error
		}
		ctx.Status(201)
		if err := ctx.JSON(ideaComment); err != nil {
			return ctx.SendStatus(500)
		}
		return nil
	}
}

func UpdateIdeaComment(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ideaComment := &models.IdeaComment{}
		db.Find(&ideaComment, ctx.Params("id"))
		if err := ctx.BodyParser(ideaComment); err != nil {
			return ctx.SendStatus(400)
		}
		result := db.Save(&ideaComment)
		if result.Error != nil {
			return result.Error
		}
		ctx.Status(200)
		if err := ctx.JSON(ideaComment); err != nil {
			return ctx.SendStatus(500)
		}
		return nil
	}
}

func DeleteIdeaComment(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ideaComment := &models.IdeaComment{}
		result := db.Delete(&ideaComment, ctx.Params("id"))
		if result.Error != nil {
			return result.Error
		}
		ctx.Status(204)
		return nil
	}
}
