package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Static("/", "./public")
	print("Starting server...")

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
