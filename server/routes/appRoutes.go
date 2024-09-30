package routes

import (
	"github.com/gofiber/fiber/v2"
	"ideashare/config"
)

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
