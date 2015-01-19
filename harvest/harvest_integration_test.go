// +build integration

package harvest_test

import (
	"os"
	"testing"

	"github.com/mitch000001/go-harvest/cmd/harvest"
	"github.com/mitch000001/go-harvest/harvest"
)

func createClient(t *testing.T) *harvest.Harvest {
	subdomain := os.Getenv("HARVEST_SUBDOMAIN")
	username := os.Getenv("HARVEST_USERNAME")
	password := os.Getenv("HARVEST_PASSWORD")

	client, err := main.NewBasicAuthClient(subdomain, &main.BasicAuthConfig{username, password})
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("Expected client not to be nil")
	}
	return client
}
