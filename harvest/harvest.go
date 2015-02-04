package harvest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
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
func NewHarvest(subdomain string, clientProvider HttpClientProvider) (*Harvest, error) {
	baseUrl, err := parseSubdomain(subdomain)
	if err != nil {
		return nil, err
	}
	logger := log.New(os.Stdout, "harvest: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Base URL: %s\n", baseUrl)
	api := &JsonApi{
		Client:  clientProvider.Client,
		baseUrl: baseUrl,
		Logger:  logger,
	}
	h := &Harvest{
		baseUrl: baseUrl,
		api:     api,
	}
	h.Users = NewUserService(api.forPath("people"))
	h.Projects = NewProjectService(api.forPath("projects"))
	h.Clients = NewClientService(api.forPath("clients"))
	h.Tasks = NewTaskService(api.forPath("tasks"), api.forPath("tasks"))
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
	fmt.Printf("Account json: '%s'\n", string(responseBytes))
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

func IsNotFound(err error) bool {
	if e, ok := err.(NotFound); ok {
		return e.NotFound()
	}
	return false
}
