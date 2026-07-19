package controllers_test

import (
	"os"
	"testing"

	"MyChatApp/monolithic/internal/testutil"
)

func TestMain(m *testing.M) {
	testutil.Setup()
	code := m.Run()
	testutil.Teardown()
	os.Exit(code)
}
