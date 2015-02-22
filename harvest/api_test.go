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
	"sort"
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
		name:           "Test",
		marshaledValue: testJson,
	}
	marshaled, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}
	t.setResponseBody(statusCode, ioutil.NopCloser(bytes.NewBuffer(marshaled)))
	for k, v := range header {
		t.testResponse.Header[k] = v
	}
}

func (t *testHttpClient) setResponsePayloadAsArray(statusCode int, data interface{}) {
	testJson, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	payload := []*JsonApiPayload{
		&JsonApiPayload{
			name:           "Test",
			marshaledValue: testJson,
		},
	}
	marshaled, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}
	t.setResponseBody(statusCode, ioutil.NopCloser(bytes.NewBuffer(marshaled)))
}

func (t *testHttpClient) setResponseBody(statusCode int, body io.ReadCloser) {
	if t.testResponse == nil {
		t.testResponse = &http.Response{}
	}
	t.testResponse.StatusCode = statusCode
	t.testResponse.Body = body
	header := make(http.Header, 0)
	header.Add("Content-Type", "application/json; charset=utf-8")
	t.testResponse.Header = header
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

func sortBytes(data []byte) []byte {
	output := make([]byte, len(data))
	copy(output, data)
	sort.Sort(sortedBytes(output))
	return output
}

type sortedBytes []byte

func (s sortedBytes) Len() int           { return len(s) }
func (s sortedBytes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortedBytes) Less(i, j int) bool { return s[i] < s[j] }

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
	CrudTogglerEndpoint
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
