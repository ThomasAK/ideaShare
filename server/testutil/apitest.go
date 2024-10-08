package testutil

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
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

type FakeVerifier struct {
}

func (v *FakeVerifier) Verify(_ context.Context, _ string) (*oidc.IDToken, error) {
	return &oidc.IDToken{
		Issuer:          "test",
		Audience:        nil,
		Subject:         "test-user",
		Expiry:          time.Time{},
		IssuedAt:        time.Time{},
		Nonce:           "1234",
		AccessTokenHash: "1234",
	}, nil
}

func SetupApiTest(afterAppSetup func(*fiber.App, *config.AppContainer)) {
	InitDB()
	App, Container = app.RunApp(config.MakeDsn("", "", "localhost", "3319", "ideashare", true))
	Container.IdTokenVerifier = &FakeVerifier{}

	afterAppSetup(App, Container)
	go func() {
		if err := App.Listen(":3031"); err != nil {
			panic(err)
		}
	}()
	WaitForHealthy(10 * time.Second)
	AdminUser = &models.User{
		ExternalID: "test-user",
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
