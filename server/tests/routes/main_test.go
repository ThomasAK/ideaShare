package routes

import (
	"ideashare/tests/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetupApiTest()
	code := m.Run()
	testutil.TeardownApiTest()
	os.Exit(code)
}
