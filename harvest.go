package main

import (
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const basePathTemplate = "https://%s.harvestapp.com/"

func parseSubdomain(subdomain string) (*url.URL, error) {
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

// new creates a new Client for the provided subdomain
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

type UsersService struct {
	h *Client
}

func NewUsersService(client *Client) *UsersService {
	return &UsersService{h: client}
}

type User struct {
	Id                           int       `json:"id"`
	Email                        string    `json:"email"`
	FirstName                    string    `json:"first-name"`
	LastName                     string    `json:"last-name"`
	HasAccessToAllFutureProjects bool      `json:"has-access-to-all-future-projects"`
	DefaultHourlyRate            int       `json:"default-hourly-rate"`
	IsActive                     bool      `json:"is-active"`
	IsAdmin                      bool      `json:"is-admin"`
	IsContractor                 bool      `json:"is-contractor"`
	Telephone                    string    `json:"telephone"`
	Department                   string    `json:"department"`
	Timezone                     string    `json:"timezone"`
	UpdatedAt                    time.Time `json:"updated-at"`
	CreatedAt                    time.Time `json:"created-at"`
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
	users := make([]*User, 0)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(users)
	if err != nil {
		return nil, err
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
	user := User{}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UsersService) Create(user *User) (*User, error) {
	marshaledUser, err := json.Marshal(user)
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
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&user)
	if err != nil {
		return nil, err
	}
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
	marshaledUser, err := json.Marshal(user)
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
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&user)
	if err != nil {
		return nil, err
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
