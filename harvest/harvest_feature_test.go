// +build feature

package harvest_test

import (
	"os"
	"testing"

	"github.com/mitch000001/go-harvest/harvest"
	"github.com/mitch000001/go-harvest/harvest/auth"
)

func createClient(t *testing.T) *harvest.Harvest {
	subdomain := os.Getenv("HARVEST_SUBDOMAIN")
	username := os.Getenv("HARVEST_USERNAME")
	password := os.Getenv("HARVEST_PASSWORD")

	config := auth.BasicAuthConfig{
		Username: username,
		Password: password,
	}

	client, err := harvest.New(subdomain, auth.NewBasicAuthClientProvider(&config).Client)
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("Expected client not to be nil")
	}
	return client
}
