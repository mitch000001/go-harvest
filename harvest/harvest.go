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
)

const basePathTemplate = "https://%s.harvestapp.com/"

// parseSubdomain parses the subdomain string and returns a fully qualifying URL.
// It returns an error if the given string is the empty string or the string
// can't be parsed as url.URL
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

// NewHarvest creates a new Client
//
// The subdomain must either be only the subdomain or the fully qualified url.
// The clientProvider is a function providing the HttpClient used by the client.
//
// It returns an error if the subdomain does not satisfy the above mentioned specification
// or if the URL parsed from the subdomain string is not valid.
func NewHarvest(subdomain string, clientProvider HttpClientProvider) (*Harvest, error) {
	baseUrl, err := parseSubdomain(subdomain)
	if err != nil {
		return nil, err
	}
	h := &Harvest{
		baseUrl: baseUrl,
		api:     &JsonApi{Client: clientProvider.Client},
	}
	h.Users = NewUsersService(h)
	h.Projects = NewProjectsService(h)
	h.Clients = NewClientsService(h)
	return h, nil
}

// Harvest defines the client for requests on the API
type Harvest struct {
	api      CrudEndpoint
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
	response, err := h.api.(*JsonApi).Client().Do(request)
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
