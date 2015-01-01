// +build integration

package harvest

import (
	"os"
	"testing"
)

func createClient(t *testing.T) *Harvest {
	subdomain := os.Getenv("HARVEST_SUBDOMAIN")
	username := os.Getenv("HARVEST_USERNAME")
	password := os.Getenv("HARVEST_PASSWORD")

	client, err := NewBasicAuthClient(subdomain, &BasicAuthConfig{username, password})
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("Expected client not to be nil")
	}
	return client
}
