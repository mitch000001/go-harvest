package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const basePath = "https://api.harvestapp.com/"

func New(client *http.Client) (*HarvestClient, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	h := &HarvestClient{client: client, BasePath: basePath}
	h.Users = NewUsersService(h)
	return h, nil
}

type HarvestClient struct {
	client   *http.Client
	BasePath string // API endpoint base URL
	Users    *UsersService
}

func (h *HarvestClient) AbsoluteUrl(relativeRequestUrl string) string {
	if strings.HasSuffix(h.BasePath, "/") {
		return h.BasePath[:len(h.BasePath)-1] + relativeRequestUrl
	} else {
		return h.BasePath + relativeRequestUrl
	}
}

type UsersService struct {
	h *HarvestClient
}

func NewUsersService(client *HarvestClient) *UsersService {
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
	urls := s.h.AbsoluteUrl(peopleUrl)
	response, err := s.h.client.Get(urls)
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
	urls := s.h.AbsoluteUrl(fmt.Sprintf("/people/%d", id))
	response, err := s.h.client.Get(urls)
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
	urls := s.h.AbsoluteUrl("/people")
	response, err := s.h.client.Post(urls, "application/json", bytes.NewBuffer(marshaledUser))
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
	urls := s.h.AbsoluteUrl(fmt.Sprintf("/people/%d/reset_password", user.Id))
	_, err = s.h.client.Post(urls, "application/json", bytes.NewBuffer(marshaledUser))
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
	urls := s.h.AbsoluteUrl(fmt.Sprintf("/people/%d", user.Id))
	request, err := http.NewRequest("PUT", urls, bytes.NewBuffer(marshaledUser))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := s.h.client.Do(request)
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
	urls := s.h.AbsoluteUrl(fmt.Sprintf("/people/%d", user.Id))
	request, err := http.NewRequest("DELETE", urls, nil)
	if err != nil {
		return false, err
	}
	response, err := s.h.client.Do(request)
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
	urls := s.h.AbsoluteUrl(fmt.Sprintf("/people/%d", user.Id))
	request, err := http.NewRequest("POST", urls, nil)
	if err != nil {
		return false, err
	}
	response, err := s.h.client.Do(request)
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
