package harvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	m := f.(map[string]interface{})
	for k := range m {
		a.name = k
	}
	val := m[a.name]
	raw, err := json.Marshal(val)
	if err != nil {
		return err
	}
	a.marshaledValue = raw
	return nil
}

type JsonApi struct {
	baseUrl *url.URL          // API base URL
	path    string            // API endpoint path
	Client  func() HttpClient // HTTP Client to do the requests
	Logger  *log.Logger
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

func (a *JsonApi) logf(format string, arg ...interface{}) {
	if a.Logger == nil {
		return
	}
	a.Logger.Printf(format, arg)
}

func (a *JsonApi) forPath(path string) *JsonApi {
	return &JsonApi{
		baseUrl: a.baseUrl,
		path:    path,
		Client:  a.Client,
		Logger:  a.Logger,
	}
}

func (a *JsonApi) Process(method string, path string, body io.Reader) (*http.Response, error) {
	requestUrl, err := a.baseUrl.Parse(path)
	if err != nil {
		a.logf("Error parsing path: %s\n", path)
		a.logf("%T: %v\n", err, err)
		return nil, err
	}
	request, err := http.NewRequest(method, requestUrl.String(), body)
	if err != nil {
		a.logf("Error creating new request: %s\n")
		a.logf("%T: %v\n", err, err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application/json; charset=utf-8")
	response, err := a.Client().Do(request)
	if err != nil {
		return nil, err
	}
	// TODO: adapt tests to always get a response if err is nil
	if ct := response.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		return nil, fmt.Errorf("Bad Request: \nResponse has wrong Content-Type '%q'\nRequest: %+#v\nRequest URL: %s\n", ct, request, request.URL)
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
		a.logf("%T: %v\n", err, err)
		return err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	var payload []*JsonApiPayload
	err = json.Unmarshal(responseBytes, &payload)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		a.logf("Response: %s\n", string(responseBytes))
		return err
	}
	var rawPayloads []*json.RawMessage
	for _, p := range payload {
		rawPayloads = append(rawPayloads, p.MarshaledValue())
	}
	marshaled, err := json.Marshal(&rawPayloads)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	err = json.Unmarshal(marshaled, &data)
	if err != nil {
		a.logf("%T: %v\n", err, err)
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
		a.logf("%T: %v\n", err, err)
		return err
	}
	if response.StatusCode == 404 {
		a.logf("%T: %v\n", err, err)
		return notFound("")
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	var payload JsonApiPayload
	err = json.Unmarshal(responseBytes, &payload)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	marshaled, err := json.Marshal(payload.MarshaledValue())
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	err = json.Unmarshal(marshaled, data)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	return nil
}

// Create creates a new data entry at the API endpoint
func (a *JsonApi) Create(data CrudModel) error {
	marshaledData, err := json.Marshal(&data)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	requestPayload := &JsonApiPayload{
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}

	response, err := a.Process("POST", a.path, bytes.NewReader(marshaledPayload))
	if err != nil {
		a.logf("%T: %v\n", err, err)
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
			a.logf("%T: %v\n", err, err)
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			a.logf("%T: %v\n", err, err)
			return err
		}
		err = &ResponseError{&apiResponse}
		a.logf("%T: %v\n", err, err)
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
		a.logf("%T: %v\n", err, err)
		return err
	}
	requestPayload := &JsonApiPayload{
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	response, err := a.Process("PUT", fmt.Sprintf(updateTemplate, id), bytes.NewReader(marshaledPayload))
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			a.logf("%T: %v\n", err, err)
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			a.logf("%T: %v\n", err, err)
			return err
		}
		err = &ResponseError{&apiResponse}
		a.logf("%T: %v\n", err, err)
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
		a.logf("%T: %v\n", err, err)
		return err
	}
	requestPayload := &JsonApiPayload{
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}

	response, err := a.Process("DELETE", fmt.Sprintf(deleteTemplate, id), bytes.NewReader(marshaledPayload))
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			a.logf("%T: %v\n", err, err)
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			a.logf("%T: %v\n", err, err)
			return err
		}
		err = &ResponseError{&apiResponse}
		a.logf("%T: %v\n", err, err)
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
		a.logf("%T: %v\n", err, err)
		return err
	}
	requestPayload := &JsonApiPayload{
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}

	response, err := a.Process("POST", fmt.Sprintf(toggleTemplate, id), bytes.NewReader(marshaledPayload))
	if err != nil {
		a.logf("%T: %v\n", err, err)
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		data.ToggleActive()
	} else if response.StatusCode == http.StatusBadRequest {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			a.logf("%T: %v\n", err, err)
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			a.logf("%T: %v\n", err, err)
			return err
		}
		err = &ResponseError{&apiResponse}
		a.logf("%T: %v\n", err, err)
		return err
	} else {
		panic(fmt.Sprintf("Unknown StatusCode: %d", response.StatusCode))
	}
	return nil
}
