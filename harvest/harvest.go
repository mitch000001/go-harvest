package harvest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"code.google.com/p/goauth2/oauth"
)

const basePathTemplate = "https://%s.harvestapp.com/"

func parseSubdomain(subdomain string) (*url.URL, error) {
	if subdomain == "" {
		return nil, errors.New("Subdomain can't be blank")
	}
	if len(strings.Split(subdomain, ".")) == 1 {
		return url.Parse(fmt.Sprintf(basePathTemplate, subdomain))
	}
	if !strings.HasSuffix(subdomain, "/") {
		subdomain = subdomain + "/"
	}
	return url.Parse(subdomain)
}

type authenticationTransport interface {
	Client() *http.Client
}

// newHarvest creates a new Client for the provided subdomain
func newHarvest(subdomain string) (*Harvest, error) {
	baseUrl, err := parseSubdomain(subdomain)
	if err != nil {
		return nil, err
	}
	h := &Harvest{baseUrl: baseUrl}
	h.Users = NewUsersService(h)
	h.Projects = NewProjectsService(h)
	h.Clients = NewClientsService(h)
	return h, nil
}

// NewBasicAuthClient creates a new Client with BasicAuth as authentication method
func NewBasicAuthClient(subdomain string, config *BasicAuthConfig) (*Harvest, error) {
	h, err := newHarvest(subdomain)
	if err != nil {
		return nil, err
	}
	h.authenticationTransport = &Transport{Config: config}
	return h, nil
}

// NewOAuthClient creates a new Client with OAuth as authentication method
func NewOAuthClient(subdomain string, config *oauth.Config) (*Harvest, error) {
	h, err := newHarvest(subdomain)
	if err != nil {
		return nil, err
	}
	h.authenticationTransport = &oauth.Transport{Config: config}
	return h, err
}

// Harvest defines the client for requests on the API
type Harvest struct {
	authenticationTransport
	baseUrl  *url.URL // API endpoint base URL
	Users    *UsersService
	Projects *ProjectsService
	Clients  *ClientsService
}

func (h *Harvest) ProcessRequest(method string, path string, body io.Reader) (*http.Response, error) {
	requestUrl, err := h.baseUrl.Parse(path)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, requestUrl.String(), body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	if err != nil {
		return nil, err
	}
	response, err := h.Client().Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (h *Harvest) Account() (*Account, error) {
	response, err := h.ProcessRequest("GET", "/account/who_am_i", nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	account := Account{}
	err = json.Unmarshal(responseBytes, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

type ErrorPayload struct {
	Message string `json:"message,omitempty"`
}

type ResponseError struct {
	ErrorPayload *ErrorPayload
}

func (r *ResponseError) Error() string {
	return r.ErrorPayload.Message
}

type NotFound interface {
	error
	NotFound() bool
}

func notFound(message string) NotFound {
	if message == "" {
		message = "Not found"
	}
	return NotFoundError(message)
}

type NotFoundError string

func (n NotFoundError) Error() string {
	return string(n)
}

func (n NotFoundError) NotFound() bool {
	return true
}

func isNotFound(err error) bool {
	if e, ok := err.(NotFound); ok {
		return e.NotFound()
	}
	return false
}

type ApiPayload struct {
	Name  string
	Value json.RawMessage
}

var apiPayloadJSONTemplate string = `{"%s": %s}`

func (a *ApiPayload) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(apiPayloadJSONTemplate, a.Name, a.Value)), nil
}

func (a *ApiPayload) UnmarshalJSON(data []byte) error {
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

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Api struct {
	baseUrl *url.URL          // API base URL
	path    string            // API endpoint path
	Client  func() HttpClient // HTTP Client to do the requests
}

func (a *Api) processRequest(method string, path string, body io.Reader) (*http.Response, error) {
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
	if err != nil {
		return nil, err
	}
	response, err := a.Client().Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// All populates the data passed in with the results found at the API endpoint
// data must be a slice of pointers to the resource corresponding with the
// endpoint
// params contains additional query parameters and may be nil
func (a *Api) All(data interface{}, params url.Values) error {
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
	var payload []*ApiPayload
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
	err = json.Unmarshal(marshaled, data)
	if err != nil {
		return err
	}
	return nil
}

// Find gets the data specified by id
// id is accepted as primitive data type or as type which implements
// the fmt.Stringer interface
func (a *Api) Find(id, data interface{}) error {
	// TODO: It's nice to build "templates" for Sprintf, but it's not comprehensible
	findTemplate := fmt.Sprintf("%s/%%%%%%c", a.path)
	idVerb := 'v'
	_, ok := id.(fmt.Stringer)
	if ok {
		idVerb = 's'
	}
	pathTemplate := fmt.Sprintf(findTemplate, idVerb)
	response, err := a.processRequest("GET", fmt.Sprintf(pathTemplate, id), nil)
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
	var payload ApiPayload
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

func (a *Api) Create(data interface{}) error {
	return errors.New("Not implemented yet")
}

func (a *Api) Update(data interface{}) error {
	return errors.New("Not implemented yet")
}

func (a *Api) Delete(data interface{}) error {
	return errors.New("Not implemented yet")
}
