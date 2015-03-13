package harvest

import (
	"fmt"
	"testing"
)

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
	testClientFn := func() HttpClient { return &testHttpClient{} }
	client, err := New("foo", testClientFn)

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

	if client.Tasks == nil {
		t.Logf("Expected tasks service not to be nil")
		t.Fail()
	}

	// wrong kind of subdomain
	client, err = New("", testClientFn)

	if err == nil {
		t.Logf("Expected error\n")
		t.Fail()
	}

	if client != nil {
		t.Logf("Expected returning client to be nil\n")
		t.Fail()
	}
}

func TestNotFound(t *testing.T) {
	notFoundError := notFound("foo")

	errMessage := notFoundError.Error()

	expectedMessage := "foo"

	if errMessage != expectedMessage {
		t.Logf("Expected message to equal 'q', got '%q'\n", expectedMessage, errMessage)
		t.Fail()
	}

	// No message given
	notFoundError = notFound("")

	errMessage = notFoundError.Error()

	expectedMessage = "Not found"

	if errMessage != expectedMessage {
		t.Logf("Expected message to equal 'q', got '%q'\n", expectedMessage, errMessage)
		t.Fail()
	}
}

func TestNotFoundNotFound(t *testing.T) {
	notFoundError := notFound("")

	ok := notFoundError.NotFound()

	if !ok {
		t.Logf("Expected NotFound to return true, got false\n")
		t.Fail()
	}
}

type found string

func (f found) Error() string {
	return string(f)
}

func (f found) NotFound() bool {
	return false
}

func TestIsNotFound(t *testing.T) {
	notFoundError := notFound("")

	ok := IsNotFound(notFoundError)

	if !ok {
		t.Logf("Expected IsNotFound to return true, got false\n")
		t.Fail()
	}

	// Any other error
	err := fmt.Errorf("foo")

	ok = IsNotFound(err)

	if ok {
		t.Logf("Expected IsNotFound to return false, got true\n")
		t.Fail()
	}

	// An error implementing NotFound
	err = found("baz")

	ok = IsNotFound(err)

	if ok {
		t.Logf("Expected IsNotFound to return false, got true\n")
		t.Fail()
	}
}
