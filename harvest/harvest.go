package harvest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

// NewHarvest creates a new Client
//
// The subdomain must either be only the subdomain or the fully qualified url.
// The clientProvider is a function providing the HttpClient used by the client.
//
// It returns an error if the subdomain does not satisfy the above mentioned specification
// or if the URL parsed from the subdomain string is not valid.
func New(subdomain string, clientProvider func() HttpClient) (*Harvest, error) {
	baseUrl, err := parseSubdomain(subdomain)
	if err != nil {
		return nil, err
	}
	api := &JsonApi{
		Client:  clientProvider,
		baseUrl: baseUrl,
	}
	h := &Harvest{
		baseUrl: baseUrl,
		api:     api,
	}
	userApi := api.CrudTogglerEndpoint("people")
	h.Users = NewUserService(api, userApi)
	projectApi := api.CrudTogglerEndpoint("projects")
	h.Projects = NewProjectService(api, projectApi)
	h.Clients = NewClientService(api.CrudTogglerEndpoint("clients"))
	taskApi := api.CrudTogglerEndpoint("tasks")
	h.Tasks = NewTaskService(taskApi, api)
	return h, nil
}

// Harvest defines the client for requests on the API
type Harvest struct {
	api      *JsonApi
	baseUrl  *url.URL // API endpoint base URL
	Users    *UserService
	Projects *ProjectService
	Clients  *ClientService
	Tasks    *TaskService
}

func (h *Harvest) Account() (*Account, error) {
	response, err := h.api.Process("GET", "/account/who_am_i", nil)
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

func (h *Harvest) RateLimitStatus() (interface{}, error) {
	response, err := h.api.Process("GET", "account/rate_limit_status", nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var limit RateLimit
	err = json.Unmarshal(responseBytes, &limit)
	if err != nil {
		return nil, err
	}
	return limit, nil
}

type RateLimit struct {
	// TimeframeLimit specifies the timframe for the quota. It is provided in seconds.
	TimeframeLimit int `json:"timeframe_limit"`
	// MaxCalls defines the maximum quota per timeframe
	MaxCalls int `json:"max_calls"`
	// Count provides the API calls
	Count             int `json:"count"`
	RequestsAvailable int `json:"requests_available"`
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

type RateLimitReached interface {
	error
	RateLimitReached() bool
	RetryAfter() string
}

func rateLimitReached(message string) RateLimitReachedError {
	if message == "" {
		message = "Rate limit reached"
	}
	return RateLimitReachedError(message)
}

type RateLimitReachedError string

func (r RateLimitReachedError) Error() string {
	return string(r)
}

func (r RateLimitReachedError) RateLimitReached() bool {
	return true
}

func (r RateLimitReachedError) Temporary() bool {
	return true
}

func IsRateLimitReached(err error) bool {
	if e, ok := err.(RateLimitReached); ok {
		return e.RateLimitReached()
	}
	return false
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

func IsNotFound(err error) bool {
	if e, ok := err.(NotFound); ok {
		return e.NotFound()
	}
	return false
}
