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

func TestJsonApiAll(t *testing.T) {
	testClient := &testHttpClient{}
	api := createJsonTestApi(testClient)

	testData := testPayload{
		ID:   12,
		Data: "foobar",
	}
	testClient.setResponsePayloadAsArray(http.StatusOK, testData)

	var data []*testPayload

	err := api.All(&data, nil)

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
	testClient.setResponseBody(http.StatusOK, emptyReadCloser())

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

func TestJsonApiFind(t *testing.T) {
	testClient := &testHttpClient{}
	api := createJsonTestApi(testClient)
	testData := testPayload{
		ID:   12,
		Data: "foobar",
	}
	testClient.setResponsePayload(http.StatusOK, nil, testData)

	var data *testPayload

	err := api.Find(12, &data, nil)

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
	testClient.setResponseBody(http.StatusNotFound, emptyReadCloser())

	data = nil

	err = api.Find(999, &data, nil)

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

	// Testing url query params
	testClient.setResponseBody(http.StatusOK, emptyReadCloser())

	data = nil
	params := url.Values{}
	params.Add("foo", "bar")

	err = api.Find(12, &data, params)

	testRequestUrl := testClient.testRequest.URL

	if !reflect.DeepEqual(testRequestUrl.Query(), params) {
		t.Logf("Expected query params from request to be '%v', got: '%v'", params, testRequestUrl.Query())
		t.Fail()
	}
}

func TestJsonApiCreate(t *testing.T) {
	testClient := &testHttpClient{}
	api := createJsonTestApi(testClient)
	testData := testPayload{
		Data: "foobar",
	}

	header := http.Header{"Location": []string{fmt.Sprintf("/%s/4", api.path)}}
	testClient.setResponsePayload(http.StatusCreated, header, nil)

	err := api.Create(&testData)

	if err != nil {
		t.Logf("Expected no error, got: %v\n", err)
		t.Fail()
	}

	if testData.ID != 4 {
		t.Logf("Expected data.id to be %d, got: %d\n", 4, testData.ID)
		t.Fail()
	}

	// test invalid data
	body := &ErrorPayload{Message: "FAIL"}
	bodyBytes := panicErr(json.Marshal(&body)).([]byte)
	response := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       bytesToReadCloser(bodyBytes),
	}
	testClient.testResponse = response

	err = api.Create(&testData)

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
}

func TestJsonApiUpdate(t *testing.T) {
	testClient := &testHttpClient{}
	api := createJsonTestApi(testClient)
	testData := testPayload{
		ID:   12,
		Data: "foobar",
	}

	testClient.setResponsePayload(http.StatusOK, nil, nil)

	err := api.Update(&testData)

	if err != nil {
		t.Logf("Expected no error, got: %v\n", err)
		t.Fail()
	}

	request := testClient.testRequest
	if request == nil {
		t.Logf("Expected request not to be nil\n")
		t.FailNow()
	}
	if request.Method != "PUT" {
		t.Logf("Expected request method to be 'PUT', got '%s'\n", request.Method)
		t.Fail()
	}
	requestBodyBytes := panicErr(ioutil.ReadAll(request.Body)).([]byte)
	expectedBytes := []byte(`{"testPayload":{"ID":12,"Data":"foobar"}}`)
	if !bytes.Equal(expectedBytes, requestBodyBytes) {
		t.Logf("Expected request body to equal '%s', got '%s'\n", string(expectedBytes), string(requestBodyBytes))
		t.Fail()
	}

	// Failing update
	body := &ErrorPayload{Message: "FAIL"}
	bodyBytes := panicErr(json.Marshal(&body)).([]byte)
	response := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       bytesToReadCloser(bodyBytes),
	}
	testClient.testResponse = response

	err = api.Update(&testData)

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
}

func TestJsonApiDelete(t *testing.T) {
	testClient := &testHttpClient{}
	api := createJsonTestApi(testClient)
	testData := testPayload{
		ID:   12,
		Data: "foobar",
	}

	testClient.setResponsePayload(http.StatusOK, nil, nil)

	err := api.Delete(&testData)

	if err != nil {
		t.Logf("Expected no error, got: %v\n", err)
		t.Fail()
	}

	request := testClient.testRequest
	if request == nil {
		t.Logf("Expected request not to be nil\n")
		t.FailNow()
	}
	if request.Method != "DELETE" {
		t.Logf("Expected request method to be 'DELETE', got '%s'\n", request.Method)
		t.Fail()
	}
	requestBodyBytes := panicErr(ioutil.ReadAll(request.Body)).([]byte)
	expectedBytes := []byte(`{"testPayload":{"ID":12,"Data":"foobar"}}`)
	if !bytes.Equal(expectedBytes, requestBodyBytes) {
		t.Logf("Expected request body to equal '%s', got '%s'\n", string(expectedBytes), string(requestBodyBytes))
		t.Fail()
	}

	// Failing delete
	body := &ErrorPayload{Message: "FAIL"}
	bodyBytes := panicErr(json.Marshal(&body)).([]byte)
	response := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       bytesToReadCloser(bodyBytes),
	}
	testClient.testResponse = response

	err = api.Delete(&testData)

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
}

type toggleableTestPayload struct {
	*testPayload
	IsActive bool
}

func (t *toggleableTestPayload) ToggleActive() bool {
	t.IsActive = !t.IsActive
	return t.IsActive
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
	expectedBytes := []byte(`{"toggleableTestPayload":{"ID":12,"Data":"foobar","IsActive":true}}`)
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

type testFunc func(*apiWrapperTestData, *bool) CrudTogglerApi

func testApiAllWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerApi {
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

func testApiFindWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerApi {
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

func testApiCreateWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerApi {
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

func testApiUpdateWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerApi {
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

func testApiDeleteWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerApi {
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

func testApiToggleWrapper(testData *apiWrapperTestData, called *bool) CrudTogglerApi {
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

func testApiAll(fn func(interface{}, url.Values) error) CrudTogglerApi {
	return &testApi{allFn: fn}
}

func testApiFind(fn func(interface{}, interface{}, url.Values) error) CrudTogglerApi {
	return &testApi{findFn: fn}
}

func testApiCreate(fn func(CrudModel) error) CrudTogglerApi {
	return &testApi{createFn: fn}
}

func testApiUpdate(fn func(CrudModel) error) CrudTogglerApi {
	return &testApi{updateFn: fn}
}

func testApiDelete(fn func(CrudModel) error) CrudTogglerApi {
	return &testApi{deleteFn: fn}
}

func testApiToggle(fn func(ActiveTogglerCrudModel) error) CrudTogglerApi {
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
