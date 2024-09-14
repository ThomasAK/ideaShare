package routes

import (
	"ideashare/app"
	"ideashare/config"
	"ideashare/tests/testutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func WaitForHealthy(maxWait time.Duration) {
	res, _ := http.Get("http://localhost:3031/health")
	start := time.Now()
	for res == nil || res.StatusCode != 200 {
		if time.Now().Sub(start) > maxWait {
			panic("Server did not become healthy")
		}
		time.Sleep(100 * time.Millisecond)
		res, _ = http.Get("http://localhost:3031/health")
	}
}

func TestMain(m *testing.M) {
	testutil.InitDB()
	a := app.RunApp(config.MakeDsn("", "", "localhost", "3319", "ideashare"))
	go func() {
		if err := a.Listen(":3031"); err != nil {
			panic(err)
		}
	}()
	WaitForHealthy(10 * time.Second)
	code := m.Run()
	testutil.TearDownDB()
	_ = a.Shutdown()
	os.Exit(code)
}
