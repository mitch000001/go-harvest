package harvest

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

type testHttpClient struct {
	testRequest  *http.Request
	testResponse *http.Response
	testError    error
}

func (t *testHttpClient) Do(request *http.Request) (*http.Response, error) {
	t.testRequest = request
	return t.testResponse, t.testError
}

func TestProcessRequest(t *testing.T) {
	requestMethod := "GET"
	bodyContent := []byte("BODY")
	body := bytes.NewBuffer(bodyContent)
	path := "foobar"
	url, _ := url.Parse("http://www.example.com")

	testClient := testHttpClient{}

	client := func() HttpClient {
		return &testClient
	}

	api := Api{
		baseUrl: url,
		path:    path,
		Client:  client,
	}

	// Test
	_, err := api.processRequest(requestMethod, path, body)

	// Expectations
	if err != nil {
		t.Logf("Expected to get no error, got '%v'", err)
		t.Fail()
	}

	testRequest := testClient.testRequest

	if testRequest == nil {
		t.Log("Expected request not to be nil")
		t.Fail()
	}

	if testRequest.Method != requestMethod {
		t.Logf("Expected request to have method '%s', got '%s'", requestMethod, testRequest.Method)
		t.Fail()
	}

	requestUrl := testRequest.URL.String()
	expectedUrl := url.String() + "/" + path

	if requestUrl != expectedUrl {
		t.Logf("Expected request to have URL '%s', got '%s'", expectedUrl, requestUrl)
		t.Fail()
	}

	expectedContentType := "application/json"
	actualContentType := testRequest.Header.Get("Content-Type")

	if actualContentType != expectedContentType {
		t.Logf("Expected request to have Content-Type header set to '%s', got '%s'", expectedContentType, actualContentType)
		t.Fail()
	}

	expectedAcceptHeader := "application/json"
	actualAcceptHeader := testRequest.Header.Get("Accept")

	if actualAcceptHeader != expectedAcceptHeader {
		t.Logf("Expected request to have Accept header set to '%s', got '%s'", expectedAcceptHeader, actualAcceptHeader)
		t.Fail()
	}

	actualBodyBytes, err := ioutil.ReadAll(testRequest.Body)
	if err != nil {
		t.Logf("Expected to get no error, got '%v'", err)
		t.Fail()
	}

	if !bytes.Equal(actualBodyBytes, bodyContent) {
		t.Logf("Expected request to have body '%s', got '%s'", string(bodyContent), string(actualBodyBytes))
		t.Fail()
	}
}

type testPayload struct {
	Id   int
	Data string
}

func TestApiPayloadMarshalJSON(t *testing.T) {
	testData := testPayload{
		Id:   12,
		Data: "foobar",
	}
	testJson, err := json.Marshal(&testData)
	if err != nil {
		t.Fail()
		t.Logf("Got error: %v\n", err)
	}

	payload := ApiPayload{
		Name:  "Test",
		Value: testJson,
	}

	marshaled, err := json.Marshal(&payload)
	if err != nil {
		t.Fail()
		t.Logf("Expected no error, got: %v", err)
	}

	expected := `{"Test":{"Id":12,"Data":"foobar"}}`

	if string(marshaled) != expected {
		t.Fail()
		t.Logf("Expected marshaled JSON to equal '%s', got '%s'", expected, string(marshaled))
	}

}

func TestApiPayloadUnmarshalJSON(t *testing.T) {
	testJson := `{"Test":{"Id":12,"Data":"foobar"}}`
	var payload ApiPayload

	err := json.Unmarshal([]byte(testJson), &payload)

	if err != nil {
		t.Fail()
		t.Logf("Expected no error, got: %v", err)
	}

	expected := `harvest.ApiPayload{Name:"Test", Value:json.RawMessage}`
	actual := fmt.Sprintf(`%T{Name:"%s", Value:%T}`, payload, payload.Name, payload.Value)

	if actual != expected {
		t.Fail()
		t.Logf("Expected unmarshaled JSON to equal '%s', got '%s'", expected, actual)
	}

	expectedValue := []byte(`{"Id":12,"Data":"foobar"}`)
	sort.Sort(sortedBytes(expectedValue))
	sort.Sort(sortedBytes(payload.Value))

	if !bytes.Equal(expectedValue, payload.Value) {
		t.Logf("Expected value to equal '%s', got '%s'", string(expectedValue), string(payload.Value))
		t.Fail()
	}
}

func TestApiAll(t *testing.T) {
	path := "foobar"
	uri, _ := url.Parse("http://www.example.com")

	testClient := &testHttpClient{}

	client := func() HttpClient {
		return testClient
	}

	api := Api{
		baseUrl: uri,
		path:    path,
		Client:  client,
	}

	testData := testPayload{
		Id:   12,
		Data: "foobar",
	}
	testJson, err := json.Marshal(&testData)
	if err != nil {
		t.Logf("Got error: %v\n", err)
		t.Fail()
	}
	payload := []*ApiPayload{
		&ApiPayload{
			Name:  "Test",
			Value: testJson,
		},
	}
	marshaled, err := json.Marshal(&payload)
	if err != nil {
		t.Logf("Expected no error, got: %v", err)
		t.Fail()
	}

	testClient.testResponse = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(marshaled)),
	}

	var data []*testPayload

	err = api.All(&data, nil)

	if err != nil {
		t.Logf("Expected no error, got: %v", err)
		t.Fail()
	}

	if len(data) != 1 {
		t.Logf("Expected one item in data, got: %d", len(data))
		t.FailNow()
	}

	if data[0] == nil {
		t.Logf("Expected first item in data not to be nil")
		t.FailNow()
	}

	if !reflect.DeepEqual(*data[0], testData) {
		t.Logf("Expected data to equal %+#v, got: %+#v", testData, *data[0])
		t.Fail()
	}

	// Testing url query params
	testClient.testResponse = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte{})),
	}

	data = nil
	params := url.Values{}
	params.Add("foo", "bar")

	err = api.All(&data, params)

	testRequestUrl := testClient.testRequest.URL

	if !reflect.DeepEqual(testRequestUrl.Query(), params) {
		t.Logf("Expected query params from request to be '%v', got: '%v'", params, testRequestUrl.Query())
		t.Fail()
	}

}

func TestApiFind(t *testing.T) {
	path := "foobar"
	uri, _ := url.Parse("http://www.example.com")

	testClient := &testHttpClient{}

	client := func() HttpClient {
		return testClient
	}

	api := Api{
		baseUrl: uri,
		path:    path,
		Client:  client,
	}

	testData := testPayload{
		Id:   12,
		Data: "foobar",
	}
	testJson, err := json.Marshal(&testData)
	if err != nil {
		t.Logf("Got error: %v\n", err)
		t.Fail()
	}
	payload := ApiPayload{
		Name:  "Test",
		Value: testJson,
	}
	marshaled, err := json.Marshal(&payload)
	if err != nil {
		t.Logf("Expected no error, got: %v", err)
		t.Fail()
	}

	testClient.testResponse = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(marshaled)),
	}

	var data *testPayload

	err = api.Find(12, &data)

	if err != nil {
		t.Logf("Expected no error, got: %v", err)
		t.Fail()
	}

	if data == nil {
		t.Logf("Expected to find one item, got nil")
		t.FailNow()
	}

	if !reflect.DeepEqual(*data, testData) {
		t.Logf("Expected data to equal %+#v, got: %+#v", testData, *data)
		t.Fail()
	}

	// Testing nonexisting id
	testClient.testResponse = &http.Response{
		StatusCode: 404,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}

	data = nil

	err = api.Find(999, &data)

	if err == nil {
		t.Logf("Expected an error, got: nil")
		t.Fail()
	}
	if err != nil {
		if !isNotFound(err) {
			t.Logf("Expected NotFound error, got: %v", err)
			t.Fail()
		}
	}
}

type sortedBytes []byte

func (s sortedBytes) Len() int           { return len(s) }
func (s sortedBytes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortedBytes) Less(i, j int) bool { return s[i] < s[j] }
