package harvest

import "testing"

func TestParseSubdomain(t *testing.T) {
	// Happy path
	fullQualifiedSubdomain := "https://foo.harvestapp.com/"

	testSubdomain(fullQualifiedSubdomain, t)
	// only the subdomain name given
	onlySubdomainName := "foo"

	testSubdomain(onlySubdomainName, t)

	fullQualifiedSubdomainWithoutTrailingSlash := "https://foo.harvestapp.com"

	testSubdomain(fullQualifiedSubdomainWithoutTrailingSlash, t)

	// Invalid subdomains
	noSubdomain := ""

	testUrl, err := parseSubdomain(noSubdomain)
	if err == nil {
		t.Logf("Expected error, got nil. Resulting testUrl: '%+#v'\n", testUrl)
		t.Fail()
	}
	if err != nil {

	}
}

func testSubdomain(subdomain string, t *testing.T) {
	testUrl, err := parseSubdomain(subdomain)
	if err != nil {
		t.Fatal(err)
	}
	if testUrl == nil {
		t.Fatal("Expected url not to be nil")
	}
	expectedUrl := "https://foo.harvestapp.com/"
	if testUrl.String() != expectedUrl {
		t.Fatalf("Expected schema to equal '%s', got '%s'", expectedUrl, testUrl)
	}
}

func TestNewHarvest(t *testing.T) {
	testClient := &testHttpClient{}
	testClientProvider := &testHttpClientProvider{testClient}
	client, err := NewHarvest("foo", testClientProvider)

	if err != nil {
		t.Logf("Expected no error, got %v\n", err)
		t.Fail()
	}

	if client == nil {
		t.Logf("Expected returning client not to be nil\n")
		t.FailNow()
	}

	if client.Users == nil {
		t.Logf("Expected users service not to be nil")
		t.Fail()
	}

	if client.Projects == nil {
		t.Logf("Expected projects service not to be nil")
		t.Fail()
	}

	if client.Clients == nil {
		t.Logf("Expected clients service not to be nil")
		t.Fail()
	}

	// wrong kind of subdomain
	client, err = NewHarvest("", testClientProvider)

	if err == nil {
		t.Logf("Expected error\n")
		t.Fail()
	}

	if client != nil {
		t.Logf("Expected returning client to be nil\n")
		t.Fail()
	}
}
