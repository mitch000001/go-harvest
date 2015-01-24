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

func TestNewJsonApiPayload(t *testing.T) {
	name := "foo"
	marshaledValue := []byte("bar")
	value := "bar"

	payload := NewJsonApiPayload(name, marshaledValue, &value)

	if payload.name != name {
		t.Logf("Expected name to equal '%q', got '%q'\n", name, payload.name)
		t.Fail()
	}

	sort.Sort(sortedBytes(marshaledValue))
	sort.Sort(sortedBytes(payload.marshaledValue))
	if !bytes.Equal(marshaledValue, payload.marshaledValue) {
		t.Logf("Expected marshaledValue to equal '%q', got '%q'\n", string(marshaledValue), string(payload.marshaledValue))
		t.Fail()
	}

	if !reflect.DeepEqual(payload.value, &value) {
		t.Logf("Expected value to equal '%+#v', got '%+#v'\n", &value, payload.value)
		t.Fail()
	}
}

func TestJsonApiPayloadName(t *testing.T) {
	name := "foo"
	marshaledValue := []byte("bar")

	payload := &JsonApiPayload{
		name:           name,
		marshaledValue: marshaledValue,
	}

	actualName := payload.Name()

	if actualName != name {
		t.Logf("Expected name to equal '%q', got '%q'\n", name, payload.name)
		t.Fail()
	}
}

func TestJsonApiPayloadMarshaledValue(t *testing.T) {
	name := "foo"
	marshaledValue := []byte("bar")

	payload := &JsonApiPayload{
		name:           name,
		marshaledValue: marshaledValue,
	}

	actualMarshaledValue := payload.MarshaledValue()

	sort.Sort(sortedBytes(marshaledValue))
	sort.Sort(sortedBytes(*actualMarshaledValue))
	if !bytes.Equal(marshaledValue, *actualMarshaledValue) {
		t.Logf("Expected marshaledValue to equal '%q', got '%q'\n", string(marshaledValue), string(*actualMarshaledValue))
		t.Fail()
	}
}

func TestJsonApiPayloadValue(t *testing.T) {
	name := "foo"
	marshaledValue := []byte("bar")
	value := "bar"

	payload := &JsonApiPayload{
		name:           name,
		marshaledValue: marshaledValue,
		value:          &value,
	}

	actualValue := payload.Value()

	if !reflect.DeepEqual(actualValue, &value) {
		t.Logf("Expected value to equal '%+#v', got '%+#v'\n", &value, actualValue)
		t.Fail()
	}

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
		name:           "Test",
		marshaledValue: testJson,
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
	actual := fmt.Sprintf(`%T{Name:"%s", Value:%T}`, payload, payload.name, payload.marshaledValue)

	if actual != expected {
		t.Fail()
		t.Logf("Expected unmarshaled JSON to equal '%s', got '%s'", expected, actual)
	}

	expectedValue := []byte(`{"ID":12,"Data":"foobar"}`)
	sort.Sort(sortedBytes(expectedValue))
	sort.Sort(sortedBytes(payload.marshaledValue))

	if !bytes.Equal(expectedValue, payload.marshaledValue) {
		t.Logf("Expected value to equal '%s', got '%s'", string(expectedValue), string(payload.marshaledValue))
		t.Fail()
	}
}

func TestJsonApiProcessRequest(t *testing.T) {
	testClient := &testHttpClient{}
	api := createJsonTestApi(testClient)

	path := "qux"
	requestMethod := "GET"
	bodyContent := []byte("BODY")
	body := bytes.NewReader(bodyContent)

	// Test
	_, err := api.Process(requestMethod, path, body)

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
		if _, ok := err.(NotFound); !ok {
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
	expectedBytes := []byte(`{"testpayload":{"ID":12,"Data":"foobar"}}`)
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
	expectedBytes := []byte(`{"testpayload":{"ID":12,"Data":"foobar"}}`)
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
