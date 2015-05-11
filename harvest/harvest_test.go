package harvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
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
		t.Logf("Expected message to equal '%q', got '%q'\n", expectedMessage, errMessage)
		t.Fail()
	}

	// No message given
	notFoundError = notFound("")

	errMessage = notFoundError.Error()

	expectedMessage = "Not found"

	if errMessage != expectedMessage {
		t.Logf("Expected message to equal '%q', got '%q'\n", expectedMessage, errMessage)
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

func TestHarvestAccount(t *testing.T) {
	testClient := &testHttpClient{}
	testAccount := &Account{
		Company: &Company{
			Name: "FOO Corp",
		},
		User: &AccountUser{
			User: &User{
				FirstName: "Max",
				LastName:  "Muster",
				CreatedAt: time.Unix(0, 0),
				UpdatedAt: time.Unix(0, 0),
			},
			Admin: true,
		},
	}
	testClient.setResponsePayload(200, http.Header{"Content-Type": []string{"application/json"}}, testAccount, "account")

	testClientFn := func() HttpClient {
		return testClient
	}

	harvest, err := New("foo", testClientFn)

	if err != nil {
		panic(err)
	}

	// Happy path
	account, err := harvest.Account()

	if err != nil {
		t.Logf("Expected no error, got %T: %v\n", err, err)
		t.Fail()
	}

	if account == nil {
		t.Logf("Expected account not to be nil\n")
		t.Fail()
	} else {
		if !reflect.DeepEqual(account.Company, testAccount.Company) {
			t.Logf("Expected account company to equal \n%+#v\n\tgot\n%+#v\n", testAccount.Company, account.Company)
			t.Fail()
		}
		if !reflect.DeepEqual(account.User, testAccount.User) {
			t.Logf("Expected account user to equal \n%+#v\n\tgot\n%+#v\n", testAccount.User, account.User)
			t.Fail()
		}
		if !reflect.DeepEqual(account, testAccount) {
			t.Logf("Expected account to equal \n%+#v\n\tgot\n%+#v\n", testAccount, account)
			t.Fail()
		}
	}

	// HTTP error
	testClient.testError = fmt.Errorf("HTTP DOWN")

	account, err = harvest.Account()

	if err == nil {
		t.Logf("Expected error, got nil\n")
		t.Fail()
	} else {
		msg := err.Error()
		if msg != "HTTP DOWN" {
			t.Logf("Expected error message to equal %q, got %q\n", msg, testClient.testError.Error())
			t.Fail()
		}
	}

	if account != nil {
		t.Logf("Expected account to be nil, got %+#v\n", account)
		t.Fail()
	}

	// Malformed response
	testClient.setResponseBody(200, strings.NewReader("error"))
	testClient.testError = nil

	account, err = harvest.Account()

	if err == nil {
		t.Logf("Expected error, got nil\n")
		t.Fail()
	} else {
		if _, ok := err.(*json.SyntaxError); !ok {
			t.Logf("Expected error of type '*json.SyntaxError', got %T\n", err)
			t.Fail()
		}
	}

	if account != nil {
		t.Logf("Expected account to be nil, got %+#v\n", account)
		t.Fail()
	}
}

func TestHarvestRateLimitStatus(t *testing.T) {
	testClient := &testHttpClient{}
	testRateLimit := &RateLimit{
		TimeframeLimit:    1,
		Count:             17,
		MaxCalls:          23,
		RequestsAvailable: 5,
	}

	marshaledLimit, err := json.Marshal(testRateLimit)
	if err != nil {
		panic(err)
	}

	testClient.setResponseBody(200, bytes.NewReader(marshaledLimit))

	testClientFn := func() HttpClient {
		return testClient
	}

	harvest, err := New("foo", testClientFn)

	if err != nil {
		panic(err)
	}

	// Happy path
	rateLimit, err := harvest.RateLimitStatus()

	if err != nil {
		t.Logf("Expected no error, got %T: %v\n", err, err)
		t.Fail()
	}

	if rateLimit == nil {
		t.Logf("Expected rateLimit not to be nil\n")
		t.Fail()
	} else {
		if !reflect.DeepEqual(rateLimit, testRateLimit) {
			t.Logf("Expected rateLimit to equal \n%+#v\n\tgot\n%+#v\n", testRateLimit, rateLimit)
			t.Fail()
		}
	}

	// HTTP error
	testClient.testError = fmt.Errorf("HTTP DOWN")

	rateLimit, err = harvest.RateLimitStatus()

	if err == nil {
		t.Logf("Expected error, got nil\n")
		t.Fail()
	} else {
		msg := err.Error()
		if msg != "HTTP DOWN" {
			t.Logf("Expected error message to equal %q, got %q\n", msg, testClient.testError.Error())
			t.Fail()
		}
	}

	if rateLimit != nil {
		t.Logf("Expected rateLimit to be nil, got %+#v\n", rateLimit)
		t.Fail()
	}

	// Malformed response
	testClient.setResponseBody(200, strings.NewReader("error"))
	testClient.testError = nil

	rateLimit, err = harvest.RateLimitStatus()

	if err == nil {
		t.Logf("Expected error, got nil\n")
		t.Fail()
	} else {
		if _, ok := err.(*json.SyntaxError); !ok {
			t.Logf("Expected error of type '*json.SyntaxError', got %T\n", err)
			t.Fail()
		}
	}

	if rateLimit != nil {
		t.Logf("Expected rateLimit to be nil, got %+#v\n", rateLimit)
		t.Fail()
	}
}

func TestNewRateLimitReached(t *testing.T) {
	// no message provided
	err := NewRateLimitReachedError("", 0)

	if err == nil {
		t.Logf("Expected err not to be nil\n")
		t.Fail()
	} else {
		expectedMessage := "Rate limit reached"
		msg := err.Error()
		if msg != expectedMessage {
			t.Logf("Expected message to equal %q, got %q\n", expectedMessage, msg)
			t.Fail()
		}
	}

	// custom message provided
	err = NewRateLimitReachedError("foobar", 0)

	if err == nil {
		t.Logf("Expected err not to be nil\n")
		t.Fail()
	} else {
		expectedMessage := "foobar"
		msg := err.Error()
		if msg != expectedMessage {
			t.Logf("Expected message to equal %q, got %q\n", expectedMessage, msg)
			t.Fail()
		}
	}

	// retry after
	err = NewRateLimitReachedError("", 15)

	retryAfter := err.RetryAfter()

	if retryAfter != 15 {
		t.Logf("Expected RetryAfter to return %d, got %d\n", 15, retryAfter)
		t.Fail()
	}
}

func TestRateLimitReachedErrorTemporary(t *testing.T) {
	err := NewRateLimitReachedError("", 0)

	if !err.Temporary() {
		t.Logf("Expected Temporary to return true, got false\n")
		t.Fail()
	}
}

func TestIsRateLimitReached(t *testing.T) {
	var err error
	// RateLimitError
	err = NewRateLimitReachedError("", 0)

	isRlrError := IsRateLimitReached(err)

	if !isRlrError {
		t.Logf("Expected true, got false\n")
		t.Fail()
	}

	// Any error
	err = fmt.Errorf("Foobar")

	isRlrError = IsRateLimitReached(err)

	if isRlrError {
		t.Logf("Expected false, got true\n")
		t.Fail()
	}
}
