package routes

import (
	"ideashare/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetupApiTest(ConfigureRoutes)
	code := m.Run()
	testutil.TeardownApiTest()
	os.Exit(code)
}
