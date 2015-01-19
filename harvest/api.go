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
)

// HttpClientProvider yields a function to provide an HttpClient.
type HttpClientProvider interface {
	// Client returns an HttpClient, which defined the minimal interface
	// of a http client usable by the harvest client to process http request
	Client() HttpClient
}

// HttpClient is the minimal interface which is used by the harvest client.
type HttpClient interface {
	// Do accepts an *http.Request and processes it
	//
	// See http.Client for a possible implementation
	Do(*http.Request) (*http.Response, error)
}

type CrudEndpoint interface {
	All(interface{}, url.Values) error
	Find(interface{}, interface{}, url.Values) error
	Create(CrudModel) error
	Update(CrudModel) error
	Delete(CrudModel) error
}

type TogglerEndpoint interface {
	Toggle(ActiveTogglerCrudModel) error
}

type CrudTogglerEndpoint interface {
	CrudEndpoint
	TogglerEndpoint
}

type ActiveToggler interface {
	// Implementations of ToggleActive should toggle their active state and
	// return the current status
	ToggleActive() bool
}

type CrudModel interface {
	Id() int
	SetId(int)
}

type ActiveTogglerCrudModel interface {
	ActiveToggler
	CrudModel
}

type JsonApiPayload struct {
	Name  string
	Value json.RawMessage
}

var apiPayloadJSONTemplate string = `{"%s": %s}`

func (a *JsonApiPayload) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(apiPayloadJSONTemplate, a.Name, a.Value)), nil
}

func (a *JsonApiPayload) UnmarshalJSON(data []byte) error {
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	m := f.(map[string]interface{})
	for k := range m {
		a.Name = k
	}
	val := m[a.Name]
	raw, err := json.Marshal(val)
	if err != nil {
		return err
	}
	a.Value = raw
	return nil
}

type JsonApi struct {
	baseUrl *url.URL          // API base URL
	path    string            // API endpoint path
	Client  func() HttpClient // HTTP Client to do the requests
}

func (a *JsonApi) processRequest(method string, path string, body io.Reader) (*http.Response, error) {
	requestUrl, err := a.baseUrl.Parse(path)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, requestUrl.String(), body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	response, err := a.Client().Do(request)
	if err != nil {
		return nil, err
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
	response, err := a.processRequest("GET", completePath, nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var payload []*JsonApiPayload
	err = json.Unmarshal(responseBytes, &payload)
	if err != nil {
		return err
	}
	var rawPayloads []*json.RawMessage
	for _, p := range payload {
		rawPayloads = append(rawPayloads, &p.Value)
	}
	marshaled, err := json.Marshal(&rawPayloads)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshaled, &data)
	if err != nil {
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
	response, err := a.processRequest("GET", completePath, nil)
	if err != nil {
		return err
	}
	if response.StatusCode == 404 {
		return notFound("")
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var payload JsonApiPayload
	err = json.Unmarshal(responseBytes, &payload)
	if err != nil {
		return err
	}
	marshaled, err := json.Marshal(&payload.Value)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshaled, data)
	if err != nil {
		return err
	}
	return nil
}

// Create creates a new data entry at the API endpoint
func (a *JsonApi) Create(data CrudModel) error {
	marshaledData, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	requestPayload := &JsonApiPayload{
		Name:  reflect.TypeOf(data).Elem().Name(),
		Value: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		return err
	}

	response, err := a.processRequest("POST", a.path, bytes.NewReader(marshaledPayload))
	if err != nil {
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
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			return err
		}
		return &ResponseError{&apiResponse}
	}
}

// Update updates the provided data at the API endpoint
func (a *JsonApi) Update(data CrudModel) error {
	id := data.Id()
	// TODO: It's nice to build "templates" for Sprintf, but it's not comprehensible
	updateTemplate := fmt.Sprintf("%s/%%d", a.path)
	marshaledData, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	requestPayload := &JsonApiPayload{
		Name:  reflect.TypeOf(data).Elem().Name(),
		Value: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		return err
	}
	response, err := a.processRequest("PUT", fmt.Sprintf(updateTemplate, id), bytes.NewReader(marshaledPayload))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			return err
		}
		return &ResponseError{&apiResponse}
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
		return err
	}
	requestPayload := &JsonApiPayload{
		Name:  reflect.TypeOf(data).Elem().Name(),
		Value: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		return err
	}

	response, err := a.processRequest("DELETE", fmt.Sprintf(deleteTemplate, id), bytes.NewReader(marshaledPayload))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			return err
		}
		return &ResponseError{&apiResponse}
	}
	return nil
}

func (a *JsonApi) Toggle(data ActiveTogglerCrudModel) error {
	id := data.Id()
	// TODO: It's nice to build "templates" for Sprintf, but it's not comprehensible
	toggleTemplate := fmt.Sprintf("%s/%%d", a.path)
	marshaledData, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	requestPayload := &JsonApiPayload{
		Name:  reflect.TypeOf(data).Elem().Name(),
		Value: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		return err
	}

	response, err := a.processRequest("POST", fmt.Sprintf(toggleTemplate, id), bytes.NewReader(marshaledPayload))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		data.ToggleActive()
	} else if response.StatusCode == http.StatusBadRequest {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			return err
		}
		return &ResponseError{&apiResponse}
	} else {
		panic(fmt.Sprintf("Unknown StatusCode: %d", response.StatusCode))
	}
	return nil
}
