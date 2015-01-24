package harvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type testPayload struct {
	ID   int
	Data string
}

func (t *testPayload) Type() string {
	return "testPayload"
}

func (t *testPayload) Id() int {
	return t.ID
}

func (t *testPayload) SetId(id int) {
	t.ID = id
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

type sortedBytes []byte

func (s sortedBytes) Len() int           { return len(s) }
func (s sortedBytes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortedBytes) Less(i, j int) bool { return s[i] < s[j] }

type toggleableTestPayload struct {
	*testPayload
	IsActive bool
}

func (t *toggleableTestPayload) ToggleActive() bool {
	t.IsActive = !t.IsActive
	return t.IsActive
}

func (t *toggleableTestPayload) Type() string {
	return "toggleableTestPayload"
}

func TestJsonApiToggle(t *testing.T) {
	testClient := &testHttpClient{}
	api := createJsonTestApi(testClient)
	testData := toggleableTestPayload{
		testPayload: &testPayload{
			ID:   12,
			Data: "foobar",
		},
		IsActive: true,
	}

	testClient.setResponsePayload(http.StatusOK, nil, nil)

	err := api.Toggle(&testData)

	if err != nil {
		t.Logf("Expected no error, got: %v\n", err)
		t.Fail()
	}

	request := testClient.testRequest
	if request == nil {
		t.Logf("Expected request not to be nil\n")
		t.FailNow()
	}
	if request.Method != "POST" {
		t.Logf("Expected request method to be 'POST', got '%s'\n", request.Method)
		t.Fail()
	}
	requestBodyBytes := panicErr(ioutil.ReadAll(request.Body)).([]byte)
	expectedBytes := []byte(`{"toggleabletestpayload":{"ID":12,"Data":"foobar","IsActive":true}}`)
	if !bytes.Equal(expectedBytes, requestBodyBytes) {
		t.Logf("Expected request body to equal '%s', got '%s'\n", string(expectedBytes), string(requestBodyBytes))
		t.Fail()
	}
	if testData.IsActive {
		t.Logf("Expected IsActive to be toggled to false, got true.\n")
		t.Fail()
	}

	// Failing toggle
	testData.IsActive = true
	body := &ErrorPayload{Message: "FAIL"}
	bodyBytes := panicErr(json.Marshal(&body)).([]byte)
	response := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       bytesToReadCloser(bodyBytes),
	}
	testClient.testResponse = response

	err = api.Toggle(&testData)

	if err == nil {
		t.Logf("Expected an error, got nil\n")
		t.Fail()
	}

	if err != nil {
		expectedMessage := "FAIL"
		errorMessage := err.Error()
		if expectedMessage != errorMessage {
			t.Logf("Expected error message '%s', got '%s'\n", expectedMessage, errorMessage)
			t.Fail()
		}
	}
	if !testData.IsActive {
		t.Logf("Expected IsActive not to be toggled to false, but was.\n")
		t.Fail()
	}
}

func createJsonTestApi(client *testHttpClient) *JsonApi {
	path := "foobar"
	uri, _ := url.Parse("http://www.example.com")
	clientFunc := func() HttpClient {
		return client
	}
	api := JsonApi{
		baseUrl: uri,
		path:    path,
		Client:  clientFunc,
	}
	return &api
}

func bytesToReadCloser(data []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(data))
}

func emptyReadCloser() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte{}))
}

type apiWrapperTestData struct {
	expectedIdType       reflect.Type
	expectedDataType     reflect.Type
	expectedParams       url.Values
	expectedErrorMessage string
	errors               bytes.Buffer
}

func (a *apiWrapperTestData) getErrors() string {
	return a.errors.String()
}

type testFunc func(*apiWrapperTestData, *bool) CrudTogglerEndpoint

func testApiAllWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerEndpoint {
	testFn := func(data interface{}, params url.Values) error {
		*called = true
		dataType := reflect.TypeOf(data)

		if !reflect.DeepEqual(dataType, testData.expectedDataType) {
			fmt.Fprintf(&testData.errors, "Expected data type '%q', got '%q'\n", testData.expectedDataType, dataType)
		}

		if !reflect.DeepEqual(testData.expectedParams, params) {
			fmt.Fprintf(&testData.errors, "Expected params to equal '%q', got '%q'\n", testData.expectedParams, params)
		}

		return fmt.Errorf(testData.expectedErrorMessage)
	}
	return testApiAll(testFn)
}

func testApiFindWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerEndpoint {
	testFn := func(id interface{}, data interface{}, params url.Values) error {
		*called = true
		dataType := reflect.TypeOf(data)
		if !reflect.DeepEqual(dataType, testData.expectedDataType) {
			fmt.Fprintf(&testData.errors, "Expected data type '%q', got '%q'\n", testData.expectedDataType, dataType)
		}

		idType := reflect.TypeOf(id)
		if !reflect.DeepEqual(idType, testData.expectedIdType) {
			fmt.Fprintf(&testData.errors, "Expected data type '%q', got '%q'\n", testData.expectedIdType, idType)
		}

		if !reflect.DeepEqual(testData.expectedParams, params) {
			fmt.Fprintf(&testData.errors, "Expected params to equal '%q', got '%q'\n", testData.expectedParams, params)
		}

		return fmt.Errorf(testData.expectedErrorMessage)
	}
	return testApiFind(testFn)
}

func testApiCreateWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerEndpoint {
	testFn := func(data CrudModel) error {
		*called = true
		dataType := reflect.TypeOf(data)
		if !reflect.DeepEqual(dataType, testData.expectedDataType) {
			fmt.Fprintf(&testData.errors, "Expected data type '%q', got '%q'\n", testData.expectedDataType, dataType)
		}

		return fmt.Errorf(testData.expectedErrorMessage)
	}
	return testApiCreate(testFn)
}

func testApiUpdateWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerEndpoint {
	testFn := func(data CrudModel) error {
		*called = true
		dataType := reflect.TypeOf(data)
		if !reflect.DeepEqual(dataType, testData.expectedDataType) {
			fmt.Fprintf(&testData.errors, "Expected data type '%q', got '%q'\n", testData.expectedDataType, dataType)
		}

		return fmt.Errorf(testData.expectedErrorMessage)
	}
	return testApiUpdate(testFn)
}

func testApiDeleteWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerEndpoint {
	testFn := func(data CrudModel) error {
		*called = true
		dataType := reflect.TypeOf(data)
		if !reflect.DeepEqual(dataType, testData.expectedDataType) {
			fmt.Fprintf(&testData.errors, "Expected data type '%q', got '%q'\n", testData.expectedDataType, dataType)
		}

		return fmt.Errorf(testData.expectedErrorMessage)
	}
	return testApiDelete(testFn)
}

func testApiToggleWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerEndpoint {
	testFn := func(data ActiveTogglerCrudModel) error {
		*called = true
		dataType := reflect.TypeOf(data)
		if !reflect.DeepEqual(dataType, testData.expectedDataType) {
			fmt.Fprintf(&testData.errors, "Expected data type '%q', got '%q'\n", testData.expectedDataType, dataType)
		}

		return fmt.Errorf(testData.expectedErrorMessage)
	}
	return testApiToggle(testFn)
}

func testApiAll(fn func(interface{}, url.Values) error) CrudTogglerEndpoint {
	return &testApi{allFn: fn}
}

func testApiFind(fn func(interface{}, interface{}, url.Values) error) CrudTogglerEndpoint {
	return &testApi{findFn: fn}
}

func testApiCreate(fn func(CrudModel) error) CrudTogglerEndpoint {
	return &testApi{createFn: fn}
}

func testApiUpdate(fn func(CrudModel) error) CrudTogglerEndpoint {
	return &testApi{updateFn: fn}
}

func testApiDelete(fn func(CrudModel) error) CrudTogglerEndpoint {
	return &testApi{deleteFn: fn}
}

func testApiToggle(fn func(ActiveTogglerCrudModel) error) CrudTogglerEndpoint {
	return &testApi{toggleFn: fn}
}

type testApi struct {
	allFn    func(interface{}, url.Values) error
	findFn   func(interface{}, interface{}, url.Values) error
	createFn func(CrudModel) error
	updateFn func(CrudModel) error
	deleteFn func(CrudModel) error
	toggleFn func(ActiveTogglerCrudModel) error
}

func (t *testApi) All(data interface{}, params url.Values) error {
	return t.allFn(data, params)
}

func (t *testApi) Find(id, data interface{}, params url.Values) error {
	return t.findFn(id, data, params)
}

func (t *testApi) Create(data CrudModel) error {
	return t.createFn(data)
}

func (t *testApi) Update(data CrudModel) error {
	return t.updateFn(data)
}

func (t *testApi) Delete(data CrudModel) error {
	return t.deleteFn(data)
}

func (t *testApi) Toggle(data ActiveTogglerCrudModel) error {
	return t.toggleFn(data)
}
