package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://twitter:twitter@localhost/twitter_test?sslmode=disable"
	}

	os.Exit(m.Run())
}
