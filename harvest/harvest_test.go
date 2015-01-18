package harvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
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
}

func TestProcessRequest(t *testing.T) {
	testClient := &testHttpClient{}
	api := createJsonTestApi(testClient)

	path := "qux"
	requestMethod := "GET"
	bodyContent := []byte("BODY")
	body := bytes.NewReader(bodyContent)

	// Test
	_, err := api.processRequest(requestMethod, path, body)

	// Expectations
	if err != nil {
		t.Logf("Expected to get no error, got '%v'", err)
		t.Fail()
	}

	expectedHeader := http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}

	testClient.testRequestFor(t, map[string]interface{}{
		"Method": requestMethod,
		"URL":    panicErr(api.baseUrl.Parse(path)),
		"Header": compare(expectedHeader, func(a, b interface{}) bool {
			for k, _ := range a.(http.Header) {
				expectedHeader := a.(http.Header).Get(k)
				actualHeader := b.(http.Header).Get(k)
				if !reflect.DeepEqual(expectedHeader, actualHeader) {
					return false
				}
			}
			return true
		}),
		"Body": compare(string(bodyContent), func(a, b interface{}) bool {
			expectedContentBytes := []byte(a.(string))
			actualBody := b.(io.Reader)
			actualBodyBytes := panicErr(ioutil.ReadAll(actualBody)).([]byte)
			return bytes.Equal(actualBodyBytes, expectedContentBytes)
		}),
	})
}

type testPayload struct {
	ID   int
	Data string
}

func (t *testPayload) Id() int {
	return t.ID
}

func (t *testPayload) SetId(id int) {
	t.ID = id
}

func TestJsonApiPayloadMarshalJSON(t *testing.T) {
	testData := testPayload{
		ID:   12,
		Data: "foobar",
	}
	testJson, err := json.Marshal(&testData)
	if err != nil {
		t.Fail()
		t.Logf("Got error: %v\n", err)
	}

	payload := JsonApiPayload{
		Name:  "Test",
		Value: testJson,
	}

	marshaled, err := json.Marshal(&payload)
	if err != nil {
		t.Fail()
		t.Logf("Expected no error, got: %v", err)
	}

	expected := `{"Test":{"ID":12,"Data":"foobar"}}`

	if string(marshaled) != expected {
		t.Fail()
		t.Logf("Expected marshaled JSON to equal '%s', got '%s'", expected, string(marshaled))
	}

}

func TestJsonApiPayloadUnmarshalJSON(t *testing.T) {
	testJson := `{"Test":{"ID":12,"Data":"foobar"}}`
	var payload JsonApiPayload

	err := json.Unmarshal([]byte(testJson), &payload)

	if err != nil {
		t.Fail()
		t.Logf("Expected no error, got: %v", err)
	}

	expected := `harvest.JsonApiPayload{Name:"Test", Value:json.RawMessage}`
	actual := fmt.Sprintf(`%T{Name:"%s", Value:%T}`, payload, payload.Name, payload.Value)

	if actual != expected {
		t.Fail()
		t.Logf("Expected unmarshaled JSON to equal '%s', got '%s'", expected, actual)
	}

	expectedValue := []byte(`{"ID":12,"Data":"foobar"}`)
	sort.Sort(sortedBytes(expectedValue))
	sort.Sort(sortedBytes(payload.Value))

	if !bytes.Equal(expectedValue, payload.Value) {
		t.Logf("Expected value to equal '%s', got '%s'", string(expectedValue), string(payload.Value))
		t.Fail()
	}
}

type testHttpClientProvider struct {
	testClient *testHttpClient
}

func (cp *testHttpClientProvider) Client() HttpClient {
	return cp.testClient
}

type testHttpClient struct {
	testRequest  *http.Request
	testResponse *http.Response
	testError    error
}

func (t *testHttpClient) Do(request *http.Request) (*http.Response, error) {
	t.testRequest = request
	return t.testResponse, t.testError
}

type compareTo interface {
	// compareTo compares the inputs with the caller
	compareTo(b interface{}) bool
}

type compareToWrapper struct {
	data      interface{}
	compareFn func(interface{}, interface{}) bool
}

func (c *compareToWrapper) compareTo(in interface{}) bool {
	return c.compareFn(c.data, in)
}

func (c *compareToWrapper) GoString() string {
	return fmt.Sprintf("%+#v", c.data)
}

func (c *compareToWrapper) String() string {
	return fmt.Sprintf("%s", c.data)
}

func compare(data interface{}, compareFn func(interface{}, interface{}) bool) compareTo {
	return &compareToWrapper{data: data, compareFn: compareFn}
}

func (t *testHttpClient) testRequestFor(tt *testing.T, testData map[string]interface{}) {
	testRequest := t.testRequest
	if testRequest == nil {
		tt.Logf("Expected request not to be nil")
		tt.Fail()
	}
	requestMap, err := structToMap(testRequest)
	if err != nil {
		tt.Logf("Expected no error, got: %v\n", err)
		tt.FailNow()
	}
	for k, v := range testData {
		reqValue := requestMap[k]
		if comp, ok := v.(compareTo); ok {
			if !comp.compareTo(reqValue) {
				tt.Logf("Expected %s to equal '%+#v', got '%+#v'\n", k, v, reqValue)
				tt.Fail()
			}
		} else {
			if !reflect.DeepEqual(reqValue, v) {
				tt.Logf("Expected %s to equal '%+#v', got '%+#v'\n", k, v, reqValue)
				tt.Fail()
			}
		}
	}
}

func structToMap(input interface{}) (map[string]interface{}, error) {
	inputValue := reflect.Indirect(reflect.ValueOf(input))
	if kind := inputValue.Kind(); kind != reflect.Struct {
		return nil, fmt.Errorf("Can't turn %v into map\n", kind)
	}
	inputType := inputValue.Type()
	output := make(map[string]interface{})
	for i := 0; i < inputValue.NumField(); i++ {
		fieldName := inputType.Field(i).Name
		output[fieldName] = inputValue.Field(i).Interface()
	}
	return output, nil
}

func (t *testHttpClient) setResponsePayload(statusCode int, header http.Header, data interface{}) {
	testJson, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	payload := &JsonApiPayload{
		Name:  "Test",
		Value: testJson,
	}
	marshaled, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}
	if t.testResponse == nil {
		t.testResponse = &http.Response{}
	}
	t.testResponse.StatusCode = statusCode
	t.testResponse.Body = ioutil.NopCloser(bytes.NewBuffer(marshaled))
	t.testResponse.Header = header
}

func (t *testHttpClient) setResponsePayloadAsArray(statusCode int, data interface{}) {
	testJson, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	payload := []*JsonApiPayload{
		&JsonApiPayload{
			Name:  "Test",
			Value: testJson,
		},
	}
	marshaled, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}
	if t.testResponse == nil {
		t.testResponse = &http.Response{}
	}
	t.testResponse.StatusCode = statusCode
	t.testResponse.Body = ioutil.NopCloser(bytes.NewBuffer(marshaled))
}

func (t *testHttpClient) setResponseBody(statusCode int, body io.ReadCloser) {
	if t.testResponse == nil {
		t.testResponse = &http.Response{}
	}
	t.testResponse.StatusCode = statusCode
	t.testResponse.Body = body
}

func panicErr(input interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return input
}

type sortedBytes []byte

func (s sortedBytes) Len() int           { return len(s) }
func (s sortedBytes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortedBytes) Less(i, j int) bool { return s[i] < s[j] }
