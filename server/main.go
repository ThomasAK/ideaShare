package main

import (
	"ideashare/app"
	"ideashare/config"
)

func main() {
	a := app.RunApp(config.MakeDsn(
		config.GetStringOr(config.DbUser, "ideashare"),
		config.GetStringOr(config.DbPass, "password"),
		config.GetStringOr(config.DbHost, "localhost"),
		config.GetStringOr(config.DbPort, "3318"),
		config.GetStringOr(config.DbName, "ideashare"),
	))
	if err := a.Listen(":3030"); err != nil {
		panic(err)
	}
}
