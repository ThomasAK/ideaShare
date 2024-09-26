package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"ideashare/config"
	"ideashare/models"
	"time"
)

var dummyUser = &models.User{
	SoftDeleteModel: models.SoftDeleteModel{
		HardDeleteModel: models.HardDeleteModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: 1,
		},
		DeletedAt: gorm.DeletedAt{},
	},
	ExternalID: "1",
	FirstName:  "Admin",
	LastName:   "Admin",
	Roles:      []*models.UserRole{{UserID: 1, Role: models.SiteAdmin}},
}

func AppRouteWithBody[T any](container *config.AppContainer, newBody func() T, handler func(container *config.AppContainer, b T, c *fiber.Ctx) (interface{}, error)) func(c *fiber.Ctx) error {
	return AppRoute(container, func(container *config.AppContainer, c *fiber.Ctx) (interface{}, error) {
		reqBody := newBody()
		if err := c.BodyParser(reqBody); err != nil {
			return nil, c.SendStatus(400)
		}
		return handler(container, reqBody, c)
	})

}

func AppRoute(container *config.AppContainer, handler func(container *config.AppContainer, c *fiber.Ctx) (interface{}, error)) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		//container.AuthProvider.authenticate(c)
		c.Locals("user", dummyUser)
		res, err := handler(container, c)
		if err != nil || res == nil {
			return err
		}
		c.Append("Content-Type", "application/json")
		if err = c.JSON(res); err != nil {
			return c.SendStatus(500)
		}
		return nil
	}
}
