package harvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var apiPayloadJSONTemplate string = `{"%s":%s}`

type JsonApiPayload struct {
	name           string
	marshaledValue json.RawMessage
	value          interface{}
}

func NewJsonApiPayload(name string, marshaledValue json.RawMessage, value interface{}) *JsonApiPayload {
	return &JsonApiPayload{
		name:           name,
		marshaledValue: marshaledValue,
		value:          value,
	}
}

func (a *JsonApiPayload) Name() string {
	return a.name
}

func (a *JsonApiPayload) MarshaledValue() *json.RawMessage {
	return &a.marshaledValue
}

func (a *JsonApiPayload) Value() interface{} {
	return a.value
}

func (a *JsonApiPayload) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(apiPayloadJSONTemplate, a.name, a.marshaledValue)), nil
}

func (a *JsonApiPayload) UnmarshalJSON(data []byte) error {
	// check for syntax
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	// TODO: proper scanning and parsing from json!
	buffer := bytes.NewBuffer(data)
	buffer.ReadBytes(byte('{')) // beginning of JSON
	buffer.ReadBytes(byte('"')) // name left quote
	nameWithQuote, _ := buffer.ReadBytes(byte('"'))
	name := nameWithQuote[:len(nameWithQuote)-1]
	buffer.ReadBytes(byte(':')) // name right quote and colon
	innerJson, err := buffer.ReadBytes(byte('}'))
	if err != nil {
		return fmt.Errorf("No inner JSON!")
	}
	a.name = string(name)
	a.marshaledValue = innerJson
	return nil
}

type JsonApi struct {
	baseUrl *url.URL          // API base URL
	path    string            // API endpoint path
	Client  func() HttpClient // HTTP Client to do the requests
}

func (a *JsonApi) URL() url.URL {
	return *a.baseUrl
}

func (a *JsonApi) Path() string {
	return a.path
}

func (a *JsonApi) CrudEndpoint(path string) CrudEndpoint {
	return a.forPath(path)
}

func (a *JsonApi) TogglerEndpoint(path string) TogglerEndpoint {
	return a.forPath(path)
}

func (a *JsonApi) CrudTogglerEndpoint(path string) CrudTogglerEndpoint {
	return a.forPath(path)
}

func (a *JsonApi) forPath(path string) *JsonApi {
	return &JsonApi{
		baseUrl: a.baseUrl,
		path:    path,
		Client:  a.Client,
	}
}

func (a *JsonApi) Process(method string, path string, body io.Reader) (*http.Response, error) {
	requestUrl, err := a.baseUrl.Parse(path)
	if err != nil {
		info.Printf("Error parsing path: %s\n", path)
		info.Printf("%T: %v\n", err, err)
		return nil, err
	}
	request, err := http.NewRequest(method, requestUrl.String(), body)
	if err != nil {
		info.Printf("Error creating new request: %s\n", requestUrl.String())
		info.Printf("%T: %v\n", err, err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application/json; charset=utf-8")
	response, err := a.Client().Do(request)
	if err != nil {
		return nil, err
	}
	// TODO: adapt tests to always get a response if err is nil
	if ct := response.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			body = []byte("NO BODY")
		}
		return nil, fmt.Errorf("Bad Request: \nResponse has wrong Content-Type '%q'\nRequest: %+#v\nRequest URL: %s\nResponse: %+#v\nBody: %s\n", ct, request, request.URL, response, string(body))
	}
	if response.StatusCode == http.StatusNotFound {
		response.Body.Close()
		reason := response.Header.Get("X-404-Reason")
		return nil, notFound(reason)
	}
	if response.StatusCode == http.StatusServiceUnavailable {
		response.Body.Close()
		retryAfter := response.Header.Get("Retry-After")
		return nil, rateLimitReached(retryAfter)
	}
	return response, nil
}

// All populates the data passed in with the results found at the API endpoint.
//
// data must be a slice of pointers to the resource corresponding with the
// endpoint
//
// params contains additional query parameters and may be nil
func (a *JsonApi) All(data interface{}, params url.Values) error {
	completePath := a.path
	if params != nil {
		completePath += "?" + params.Encode()
	}
	response, err := a.Process("GET", completePath, nil)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	var payload []*JsonApiPayload
	err = json.Unmarshal(responseBytes, &payload)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		debug.Printf("Response: %+#v\n", response)
		if response.Request != nil {
			debug.Printf("Request: %+#v\n", response.Request)
			debug.Printf("Request URL: %s\n", response.Request.URL.String())
		}
		debug.Printf("Response Body: %s\n", string(responseBytes))
		return err
	}
	var rawPayloads []*json.RawMessage
	for _, p := range payload {
		rawPayloads = append(rawPayloads, p.MarshaledValue())
	}
	marshaled, err := json.Marshal(&rawPayloads)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	err = json.Unmarshal(marshaled, &data)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	return nil
}

// Find gets the data specified by id.
//
// id is accepted as primitive data type or as type which implements
// the fmt.Stringer interface.
func (a *JsonApi) Find(id interface{}, data interface{}, params url.Values) error {
	// TODO: It's nice to build "templates" for Sprintf, but it's not comprehensible
	findTemplate := fmt.Sprintf("%s/%%%%%%c", a.path)
	idVerb := 'v'
	_, ok := id.(fmt.Stringer)
	if ok {
		idVerb = 's'
	}
	pathTemplate := fmt.Sprintf(findTemplate, idVerb)
	completePath := fmt.Sprintf(pathTemplate, id)
	if params != nil {
		completePath += "?" + params.Encode()
	}
	response, err := a.Process("GET", completePath, nil)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	var payload JsonApiPayload
	err = json.Unmarshal(responseBytes, &payload)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	marshaled, err := json.Marshal(payload.MarshaledValue())
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	err = json.Unmarshal(marshaled, data)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	return nil
}

// Create creates a new data entry at the API endpoint
func (a *JsonApi) Create(data CrudModel) error {
	marshaledData, err := json.Marshal(&data)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	requestPayload := &JsonApiPayload{
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}

	response, err := a.Process("POST", a.path, bytes.NewReader(marshaledPayload))
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	defer response.Body.Close()
	id := -1
	if response.StatusCode == 201 {
		location := response.Header.Get("Location")
		scanTemplate := fmt.Sprintf("/%s/%%d", a.path)
		fmt.Sscanf(location, scanTemplate, &id)
		if id == -1 {
			return fmt.Errorf("Bad request!")
		}
		data.SetId(id)
		return nil
	} else {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			info.Printf("%T: %v\n", err, err)
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			info.Printf("%T: %v\n", err, err)
			return err
		}
		err = &ResponseError{&apiResponse}
		info.Printf("%T: %v\n", err, err)
		return err
	}
}

// Update updates the provided data at the API endpoint
func (a *JsonApi) Update(data CrudModel) error {
	id := data.Id()
	// TODO: It's nice to build "templates" for Sprintf, but it's not comprehensible
	updateTemplate := fmt.Sprintf("%s/%%d", a.path)
	marshaledData, err := json.Marshal(&data)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	requestPayload := &JsonApiPayload{
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	response, err := a.Process("PUT", fmt.Sprintf(updateTemplate, id), bytes.NewReader(marshaledPayload))
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			info.Printf("%T: %v\n", err, err)
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			info.Printf("%T: %v\n", err, err)
			return err
		}
		err = &ResponseError{&apiResponse}
		info.Printf("%T: %v\n", err, err)
		return err
	}
	return nil
}

// Delete deletes the provided data at the API endpoint
func (a *JsonApi) Delete(data CrudModel) error {
	id := data.Id()
	// TODO: It's nice to build "templates" for Sprintf, but it's not comprehensible
	deleteTemplate := fmt.Sprintf("%s/%%d", a.path)
	marshaledData, err := json.Marshal(&data)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	requestPayload := &JsonApiPayload{
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}

	response, err := a.Process("DELETE", fmt.Sprintf(deleteTemplate, id), bytes.NewReader(marshaledPayload))
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			info.Printf("%T: %v\n", err, err)
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			info.Printf("%T: %v\n", err, err)
			return err
		}
		err = &ResponseError{&apiResponse}
		info.Printf("%T: %v\n", err, err)
		return err
	}
	return nil
}

func (a *JsonApi) Toggle(data ActiveTogglerCrudModel) error {
	id := data.Id()
	// TODO: It's nice to build "templates" for Sprintf, but it's not comprehensible
	toggleTemplate := fmt.Sprintf("%s/%%d", a.path)
	marshaledData, err := json.Marshal(&data)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	requestPayload := &JsonApiPayload{
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}

	response, err := a.Process("POST", fmt.Sprintf(toggleTemplate, id), bytes.NewReader(marshaledPayload))
	if err != nil {
		info.Printf("%T: %v\n", err, err)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		data.ToggleActive()
	} else if response.StatusCode == http.StatusBadRequest {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			info.Printf("%T: %v\n", err, err)
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			info.Printf("%T: %v\n", err, err)
			return err
		}
		err = &ResponseError{&apiResponse}
		info.Printf("%T: %v\n", err, err)
		return err
	} else {
		panic(fmt.Sprintf("Unknown StatusCode: %d", response.StatusCode))
	}
	return nil
}
