package main

import (
	"os"
	"testing"
)

func TestParseSubdomain(t *testing.T) {
	// only the subdomain name given
	subdomain := "foo"

	testSubdomain(subdomain, t)

	subdomain = "https://foo.harvestapp.com/"

	testSubdomain(subdomain, t)
}

func testSubdomain(subdomain string, t *testing.T) {
	testUrl, err := parseSubdomain(subdomain)
	if err != nil {
		t.Fatal(err)
	}
	if testUrl == nil {
		t.Fatal("Expected url not to be nil")
	}
	if testUrl.String() != "https://foo.harvestapp.com/" {
		t.Fatalf("Expected schema to equal 'https://foo.harvestapp.com/', got '%s'", testUrl)
	}
}

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

func testAllFunc(testFunc interface{}, t *testing.T) {
	client := createClient(t)
	users, err := client.Users.All()
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if len(users) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(users))
	}
	if users[0] == nil {
		t.Fatal("Expected user not to be nil")
	}
	for _, u := range users {
		t.Logf("User: %+#v\n", u)
	}
}
