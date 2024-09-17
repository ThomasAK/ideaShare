package main

import (
	"ideashare/app"
	"ideashare/config"
	"ideashare/routes"
)

func main() {
	a, container := app.RunApp(config.MakeDsn(
		config.GetStringOr(config.DbUser, "ideashare"),
		config.GetStringOr(config.DbPass, "password"),
		config.GetStringOr(config.DbHost, "localhost"),
		config.GetStringOr(config.DbPort, "3318"),
		config.GetStringOr(config.DbName, "ideashare"),
		false,
	))
	routes.ConfigureRoutes(a, container)
	print("Starting server...")

	if err := a.Listen(":3030"); err != nil {
		panic(err)
	}
}
