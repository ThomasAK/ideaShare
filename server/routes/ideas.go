package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"ideashare/models"
	"strconv"
)

func ListIdeas(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		size, err := strconv.Atoi(ctx.Query("size"))
		if err != nil {
			return err
		}
		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil {
			return err
		}
		var ideas []models.Idea
		result := db.Offset((page - 1) * size).Limit(size).Find(&ideas)
		if result.Error != nil {
			return result.Error
		}
		ctx.Status(200)
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
		var idea models.Idea
		result := db.Find(&idea, ctx.Params("id"))
		if result.Error != nil {
			return result.Error
		}
		if idea.ID == 0 {
			ctx.Status(404)
			return nil
		}
		ctx.Status(200)
		if err := ctx.JSON(idea); err != nil {
			return ctx.SendStatus(500)
		}
		return nil
	}
}

func UpdateIdea(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		idea := &models.Idea{}
		db.Find(&idea, ctx.Params("id"))
		if err := ctx.BodyParser(idea); err != nil {
			return ctx.SendStatus(400)
		}
		result := db.Save(&idea)
		if result.Error != nil {
			return result.Error
		}
		ctx.Status(200)
		if err := ctx.JSON(idea); err != nil {
			return ctx.SendStatus(500)
		}
		return nil
	}
}

func DeleteIdea(db *gorm.DB) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		idea := &models.Idea{}
		result := db.Delete(&idea, ctx.Params("id"))
		if result.Error != nil {
			return result.Error
		}
		ctx.Status(204)
		return nil
	}
}
