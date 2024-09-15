package testutil

import (
	"github.com/gofiber/fiber/v2"
	"ideashare/app"
	"ideashare/config"
	"ideashare/models"
	"net/http"
	"time"
)

var (
	Container *config.AppContainer
	App       *fiber.App
	AdminUser *models.User
)

func TestApiUrl(path string) string {
	return "http://localhost:3031/api" + path
}

func WaitForHealthy(maxWait time.Duration) {
	url := "http://localhost:3031/health"
	res, _ := http.Get(url)
	start := time.Now()
	for res == nil || res.StatusCode != 200 {
		if time.Now().Sub(start) > maxWait {
			panic("Server did not become healthy")
		}
		time.Sleep(100 * time.Millisecond)
		res, _ = http.Get(url)
	}
}

func SetupApiTest() {
	InitDB()
	App, Container = app.RunApp(config.MakeDsn("", "", "localhost", "3319", "ideashare", true))
	go func() {
		if err := App.Listen(":3031"); err != nil {
			panic(err)
		}
	}()
	WaitForHealthy(10 * time.Second)
	AdminUser = &models.User{
		ExternalID: "1",
		FirstName:  "Admin",
		LastName:   "Admin",
		Roles:      []*models.UserRole{{Role: models.SiteAdmin}},
	}
	AdminUser.CreatedBy = 1
	Container.Db.Create(AdminUser)
}

func TeardownApiTest() {
	_ = App.Shutdown()
	TearDownDB()
}
