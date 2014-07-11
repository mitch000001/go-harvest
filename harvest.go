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
	"strings"
	"time"
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

// newClient creates a new Client for the provided subdomain
func newClient(subdomain string) (*Client, error) {
	baseUrl, err := parseSubdomain(subdomain)
	if err != nil {
		return nil, err
	}
	h := &Client{baseUrl: baseUrl}
	h.Users = NewUsersService(h)
	return h, nil
}

// NewBasicAuthClient creates a new Client with BasicAuth as authentication method
func NewBasicAuthClient(subdomain string, config *BasicAuthConfig) (*Client, error) {
	h, err := newClient(subdomain)
	if err != nil {
		return nil, err
	}
	h.authenticationTransport = &Transport{Config: config}
	return h, nil
}

// NewOAuthClient creates a new Client with OAuth as authentication method
func NewOAuthClient(subdomain string, config *oauth.Config) (*Client, error) {
	h, err := newClient(subdomain)
	if err != nil {
		return nil, err
	}
	h.authenticationTransport = &oauth.Transport{Config: config}
	return h, err
}

type Client struct {
	authenticationTransport
	baseUrl *url.URL // API endpoint base URL
	Users   *UsersService
}

func (h *Client) CreateRequest(method string, relativeUrl string, body io.Reader) (*http.Request, error) {
	requestUrl, err := h.baseUrl.Parse(relativeUrl)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, requestUrl.String(), body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	return request, nil
}

type Response struct {
	Message string `json:"message,omitempty"`
}

type ResponseError struct {
	Response Response
}

func (r *ResponseError) Error() string {
	return r.Response.Message
}

type UsersService struct {
	h *Client
}

func NewUsersService(client *Client) *UsersService {
	return &UsersService{h: client}
}

type User struct {
	Id                           int       `json:"id"`
	Email                        string    `json:"email"`
	FirstName                    string    `json:"first_name"`
	LastName                     string    `json:"last_name"`
	HasAccessToAllFutureProjects bool      `json:"has_access_to_all_future_projects"`
	DefaultHourlyRate            int       `json:"default_hourly_rate"`
	IsActive                     bool      `json:"is_active"`
	IsAdmin                      bool      `json:"is_admin"`
	IsContractor                 bool      `json:"is_contractor"`
	Telephone                    string    `json:"telephone"`
	Department                   string    `json:"department"`
	Timezone                     string    `json:"timezone"`
	UpdatedAt                    time.Time `json:"updated_at"`
	CreatedAt                    time.Time `json:"created_at"`
}

type UserRequest UserResponse

type UserResponse struct {
	Response
	User *User `json:"user"`
}

func (s *UsersService) All() ([]*User, error) {
	return s.AllUpdatedSince(time.Time{})
}

func (s *UsersService) AllUpdatedSince(updatedSince time.Time) ([]*User, error) {
	peopleUrl := "/people"
	if !updatedSince.IsZero() {
		values := url.Values{"updated-since": {updatedSince.UTC().String()}}
		peopleUrl = peopleUrl + "?" + values.Encode()
	}
	request, err := s.h.CreateRequest("GET", peopleUrl, nil)
	if err != nil {
		return nil, err
	}
	response, err := s.h.Client().Do(request)
	if err != nil {
		return nil, err
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	userResponses := make([]*UserResponse, 0)
	err = json.Unmarshal(responseBytes, &userResponses)
	if err != nil {
		return nil, err
	}
	users := make([]*User, len(userResponses))
	for i, u := range userResponses {
		users[i] = u.User
	}
	return users, nil
}

func (s *UsersService) Find(id int) (*User, error) {
	request, err := s.h.CreateRequest("GET", fmt.Sprintf("/people/%d", id), nil)
	if err != nil {
		return nil, err
	}
	response, err := s.h.Client().Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		return nil, &ResponseError{Response{fmt.Sprintf("No user found with id %d", id)}}
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	userResponse := UserResponse{}
	err = json.Unmarshal(responseBytes, &userResponse)
	if err != nil {
		return nil, err
	}
	return userResponse.User, nil
}

func (s *UsersService) Create(user *User) (*User, error) {
	marshaledUser, err := json.Marshal(&UserRequest{User: user})
	if err != nil {
		return nil, err
	}
	request, err := s.h.CreateRequest("POST", "/people", bytes.NewBuffer(marshaledUser))
	if err != nil {
		return nil, err
	}
	response, err := s.h.Client().Do(request)
	if err != nil {
		return nil, err
	}
	location := response.Header.Get("Location")
	userId := -1
	fmt.Sscanf(location, "/people/%d", &userId)
	if userId == -1 {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		apiResponse := Response{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			return nil, err
		}
		return nil, &ResponseError{apiResponse}
	}
	user.Id = userId
	return user, nil
}

func (s *UsersService) ResetPassword(user *User) error {
	marshaledUser, err := json.Marshal(user)
	if err != nil {
		return err
	}
	request, err := s.h.CreateRequest("POST", fmt.Sprintf("/people/%d/reset_password", user.Id), bytes.NewBuffer(marshaledUser))
	if err != nil {
		return err
	}
	_, err = s.h.Client().Do(request)
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersService) Update(user *User) (*User, error) {
	marshaledUser, err := json.Marshal(&UserRequest{User: user})
	if err != nil {
		return nil, err
	}
	request, err := s.h.CreateRequest("PUT", fmt.Sprintf("/people/%d", user.Id), bytes.NewBuffer(marshaledUser))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := s.h.Client().Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		apiResponse := Response{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			return nil, err
		}
		return nil, &ResponseError{apiResponse}
	}
	return user, nil
}

func (s *UsersService) Delete(user *User) (bool, error) {
	request, err := s.h.CreateRequest("DELETE", fmt.Sprintf("/people/%d", user.Id), nil)
	if err != nil {
		return false, err
	}
	response, err := s.h.Client().Do(request)
	if err != nil {
		return false, err
	}
	if response.StatusCode == 200 {
		return true, nil
	} else if response.StatusCode == 400 {
		return false, nil
	} else {
		panic(response.Status)
	}
}

func (s *UsersService) Toggle(user *User) (bool, error) {
	request, err := s.h.CreateRequest("POST", fmt.Sprintf("/people/%d", user.Id), nil)
	if err != nil {
		return false, err
	}
	response, err := s.h.Client().Do(request)
	if err != nil {
		return false, err
	}
	if response.StatusCode == 200 {
		return true, nil
	} else if response.StatusCode == 400 {
		return false, nil
	} else {
		panic(response.Status)
	}
}
