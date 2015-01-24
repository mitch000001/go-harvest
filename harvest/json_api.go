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
}

func NewJsonApiPayload(name string, marshaledValue json.RawMessage) *JsonApiPayload {
	return &JsonApiPayload{
		name:           name,
		marshaledValue: marshaledValue,
	}
}

func (a *JsonApiPayload) Name() string {
	return a.name
}

func (a *JsonApiPayload) MarshaledValue() *json.RawMessage {
	return &a.marshaledValue
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
}

func (a *JsonApi) Process(method string, path string, body io.Reader) (*http.Response, error) {
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
	response, err := a.Process("GET", completePath, nil)
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
		rawPayloads = append(rawPayloads, p.MarshaledValue())
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
	response, err := a.Process("GET", completePath, nil)
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
	marshaled, err := json.Marshal(payload.MarshaledValue())
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
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		return err
	}

	response, err := a.Process("POST", a.path, bytes.NewReader(marshaledPayload))
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
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		return err
	}
	response, err := a.Process("PUT", fmt.Sprintf(updateTemplate, id), bytes.NewReader(marshaledPayload))
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
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		return err
	}

	response, err := a.Process("DELETE", fmt.Sprintf(deleteTemplate, id), bytes.NewReader(marshaledPayload))
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
		name:           strings.ToLower(data.Type()),
		marshaledValue: marshaledData,
	}
	marshaledPayload, err := json.Marshal(&requestPayload)
	if err != nil {
		return err
	}

	response, err := a.Process("POST", fmt.Sprintf(toggleTemplate, id), bytes.NewReader(marshaledPayload))
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
