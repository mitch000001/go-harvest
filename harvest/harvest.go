package main

import (
	"bytes"
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

type AllFuncOptions interface {
	BuildUrl(string) string
}

type CrudService interface {
	AllFunc() interface{}
	FindFunc() interface{}
	CreateFunc() interface{}
	UpdateFunc() interface{}
	DeleteFunc() interface{}
	ResourcePath() string
}

func MakeCrudFunctions(service CrudService, client *Harvest, payloadType interface{}) {
	resourcePath := service.ResourcePath()
	parameterizedResourcePath := fmt.Sprintf("%s/%%d", resourcePath)
	MakeAllFunc(service.AllFunc(), client, service.ResourcePath(), payloadType)
	MakeFindFunc(service.FindFunc(), client, parameterizedResourcePath, payloadType)
	MakeCreateFunc(service.CreateFunc(), client, resourcePath, payloadType)
	MakeUpdateFunc(service.UpdateFunc(), client, parameterizedResourcePath, payloadType)
	MakeDeleteFunc(service.DeleteFunc(), client, parameterizedResourcePath, payloadType)
}

func MakeAllFunc(allFuncPointer interface{}, h *Harvest, path string, payloadType interface{}) {
	allFunction := reflect.ValueOf(allFuncPointer).Elem()
	allFnType := allFunction.Type()
	returnType := allFnType.Out(0)
	zeroReturnValue := reflect.Zero(returnType)
	allFuncBody := func(args []reflect.Value) []reflect.Value {
		pathVal := reflect.ValueOf(&path)
		if !args[0].IsNil() {
			urlVal := args[0].Elem().MethodByName("BuildUrl").Call([]reflect.Value{reflect.ValueOf(path)})[0]
			pathVal.Elem().Set(urlVal)
		}

		response, err := h.ProcessRequest("GET", pathVal.Elem().String(), nil)
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		payloadTyp := reflect.TypeOf(payloadType)
		payload := reflect.New(reflect.SliceOf(payloadTyp)).Interface()
		err = json.Unmarshal(responseBytes, &payload)
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		val := reflect.MakeSlice(returnType, 0, 0)
		payloadValue := reflect.Indirect(reflect.ValueOf(payload))
		for i := 0; i < payloadValue.Len(); i++ {
			val = reflect.Append(val, payloadValue.Index(i).Elem().Field(1))
		}
		return []reflect.Value{
			val,
			reflect.Zero(reflect.TypeOf(new(error)).Elem()),
		}
	}
	v := reflect.MakeFunc(allFnType, allFuncBody)
	allFunction.Set(v)
}

func MakeFindFunc(findFuncPointer interface{}, h *Harvest, path string, payloadType interface{}) {
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
				reflect.ValueOf(&err).Elem(),
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
				reflect.ValueOf(&err).Elem(),
			}
		}
		payload := reflect.New(reflect.TypeOf(payloadType)).Interface()
		err = json.Unmarshal(responseBytes, &payload)
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		return []reflect.Value{
			reflect.Indirect(reflect.ValueOf(payload)).Elem().Field(1),
			reflect.Zero(reflect.TypeOf(new(error)).Elem()),
		}
	}
	funcValue := reflect.MakeFunc(findFnType, findFuncBody)
	findFunction.Set(funcValue)
}

func MakeCreateFunc(createFuncPointer interface{}, h *Harvest, path string, payloadType interface{}) {
	createFunction := reflect.ValueOf(createFuncPointer).Elem()
	createFnType := createFunction.Type()
	returnType := createFnType.Out(0)
	zeroReturnValue := reflect.Zero(returnType)
	createFuncBody := func(args []reflect.Value) []reflect.Value {
		requestPayloadValue := reflect.New(reflect.TypeOf(payloadType).Elem())
		requestPayloadValue.Elem().Field(1).Set(args[0])
		marshaledUser, err := json.Marshal(requestPayloadValue.Interface())
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		response, err := h.ProcessRequest("POST", path, bytes.NewReader(marshaledUser))
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		location := response.Header.Get("Location")
		userId := -1
		pathFormat := fmt.Sprintf("%s/%%d", path)
		fmt.Sscanf(location, pathFormat, &userId)
		if userId == -1 {
			responseBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return []reflect.Value{
					zeroReturnValue,
					reflect.ValueOf(&err).Elem(),
				}
			}
			apiResponse := ErrorPayload{}
			err = json.Unmarshal(responseBytes, &apiResponse)
			if err != nil {
				return []reflect.Value{
					zeroReturnValue,
					reflect.ValueOf(&err).Elem(),
				}
			}
			err = &ResponseError{&apiResponse}
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		args[0].Elem().FieldByName("Id").SetInt(int64(userId))
		return []reflect.Value{
			args[0],
			reflect.Zero(reflect.TypeOf(new(error)).Elem()),
		}
	}
	funcValue := reflect.MakeFunc(createFnType, createFuncBody)
	createFunction.Set(funcValue)
}

func MakeUpdateFunc(updateFuncPointer interface{}, h *Harvest, path string, payloadType interface{}) {
	updateFunction := reflect.ValueOf(updateFuncPointer).Elem()
	updateFnType := updateFunction.Type()
	returnType := updateFnType.Out(0)
	zeroReturnValue := reflect.Zero(returnType)
	updateFuncBody := func(args []reflect.Value) []reflect.Value {
		requestPayloadValue := reflect.New(reflect.TypeOf(payloadType).Elem())
		requestPayloadValue.Elem().Field(1).Set(args[0])
		marshaledPayload, err := json.Marshal(requestPayloadValue.Interface())
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		id := args[0].Elem().FieldByName("Id").Int()
		formattedPath := fmt.Sprintf(path, id)
		response, err := h.ProcessRequest("PUT", formattedPath, bytes.NewReader(marshaledPayload))
		if err != nil {
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		if response.StatusCode != 200 {
			responseBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return []reflect.Value{
					zeroReturnValue,
					reflect.ValueOf(&err).Elem(),
				}
			}
			apiResponse := ErrorPayload{}
			err = json.Unmarshal(responseBytes, &apiResponse)
			if err != nil {
				return []reflect.Value{
					zeroReturnValue,
					reflect.ValueOf(&err).Elem(),
				}
			}
			err = &ResponseError{&apiResponse}
			return []reflect.Value{
				zeroReturnValue,
				reflect.ValueOf(&err).Elem(),
			}
		}
		return []reflect.Value{
			args[0],
			reflect.Zero(reflect.TypeOf(new(error)).Elem()),
		}
	}
	funcValue := reflect.MakeFunc(updateFnType, updateFuncBody)
	updateFunction.Set(funcValue)
}

func MakeDeleteFunc(deleteFuncPointer interface{}, h *Harvest, path string, payloadType interface{}) {
	deleteFunction := reflect.ValueOf(deleteFuncPointer).Elem()
	deleteFnType := deleteFunction.Type()
	zeroReturnValue := reflect.Zero(reflect.TypeOf(new(error)).Elem())
	deleteFuncBody := func(args []reflect.Value) []reflect.Value {
		id := args[0].Elem().FieldByName("Id").Int()
		formattedPath := fmt.Sprintf(path, id)

		response, err := h.ProcessRequest("DELETE", formattedPath, nil)
		if err != nil {
			return []reflect.Value{
				reflect.ValueOf(&err).Elem(),
			}
		}
		if response.StatusCode == 200 {
			return []reflect.Value{
				zeroReturnValue,
			}
		} else if response.StatusCode == 400 {
			err = &ResponseError{&ErrorPayload{fmt.Sprintf("Entity not deleted: %d", id)}}
			return []reflect.Value{
				reflect.ValueOf(&err).Elem(),
			}
		} else {
			panic(response.Status)
		}
	}
	funcValue := reflect.MakeFunc(deleteFnType, deleteFuncBody)
	deleteFunction.Set(funcValue)
}
