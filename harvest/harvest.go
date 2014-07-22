package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
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

func main() {
	uriBase, _ := parseSubdomain("foo")
	uri, err := uriBase.Parse("/people")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Base: %s\n", uriBase)
	fmt.Printf("Parsed: %s\n", uri)
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

func AllTemplate(path string, payload interface{}) {}

func MakeAllFunc(allFuncPointer interface{}, h *Harvest, path string) {
	allFunction := reflect.ValueOf(allFuncPointer).Elem()
	allFnType := allFunction.Type()
	returnType := allFnType.Out(0)
	zeroReturnValue := reflect.Zero(returnType)
	allFuncBody := func([]reflect.Value) []reflect.Value {
		response, err := h.ProcessRequest("GET", path, nil)
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(err),
			}
		}
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(err),
			}
		}
		payload := PayloadForType(returnType)
		err = json.Unmarshal(responseBytes, payload)
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(err),
			}
		}
		val := reflect.MakeSlice(returnType, 0, 0)
		payloadValue := reflect.ValueOf(payload)
		for i := 0; i < payloadValue.Elem().Len(); i++ {
			val = reflect.Append(val, payloadValue.Elem().Index(i).Field(1))
		}
		return []reflect.Value{
			val,
			reflect.Zero(reflect.TypeOf(new(error)).Elem()),
		}
	}
	v := reflect.MakeFunc(allFnType, allFuncBody)
	allFunction.Set(v)
}

func MakeFindFunc(findFuncPointer interface{}, h *Harvest, path string) {
	findFunction := reflect.ValueOf(findFuncPointer).Elem()
	findFnType := findFunction.Type()
	returnType := findFnType.Out(0)
	zeroReturnValue := reflect.Zero(returnType)
	findFuncBody := func(args []reflect.Value) []reflect.Value {
		idVal := args[0]
		response, err := h.ProcessRequest("GET", fmt.Sprintf(path, idVal.Int()), nil)
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(err),
			}
		}
		if response.StatusCode == 404 {
			err = &ResponseError{&ErrorPayload{fmt.Sprintf("Id not found: %d", idVal.Int())}}
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(err),
			}
		}
		payload := PayloadForType(returnType)
		err = json.Unmarshal(responseBytes, payload)
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(err),
			}
		}
		return []reflect.Value{
			reflect.ValueOf(payload).Elem().Field(1),
			reflect.Zero(reflect.TypeOf(new(error)).Elem()),
		}
	}
	funcValue := reflect.MakeFunc(findFnType, findFuncBody)
	findFunction.Set(funcValue)
}

func PayloadForType(domainType reflect.Type) interface{} {
	switch domainType {
	case reflect.TypeOf(&Project{}):
		return &ProjectPayload{}
	case reflect.TypeOf([]*Project{}):
		return &[]ProjectPayload{}
	case reflect.TypeOf(&User{}):
		return &UserPayload{}
	case reflect.TypeOf([]*User{}):
		return &[]UserPayload{}
	default:
		// FIXME: proper error handling
		return nil
	}
}
