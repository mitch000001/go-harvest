package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

type UsersService struct {
	h    *Harvest
	Find func(int) (*User, error)
	All  func() ([]*User, error)
}

func NewUsersService(client *Harvest) *UsersService {
	service := UsersService{h: client}
	MakeFindFunc(&service.Find, client, "/people/%d")
	MakeAllFunc(&service.All, client, "/people")
	return &service
}

type User struct {
	Id                           int       `json:"id,omitempty"`
	Email                        string    `json:"email,omitempty"`
	FirstName                    string    `json:"first_name,omitempty"`
	LastName                     string    `json:"last_name,omitempty"`
	HasAccessToAllFutureProjects bool      `json:"has_access_to_all_future_projects,omitempty"`
	DefaultHourlyRate            float64   `json:"default_hourly_rate,omitempty"`
	IsActive                     bool      `json:"is_active,omitempty"`
	IsAdmin                      bool      `json:"is_admin,omitempty"`
	IsContractor                 bool      `json:"is_contractor,omitempty"`
	Telephone                    string    `json:"telephone,omitempty"`
	Department                   string    `json:"department,omitempty"`
	Timezone                     string    `json:"timezone,omitempty"`
	UpdatedAt                    time.Time `json:"updated_at,omitempty"`
	CreatedAt                    time.Time `json:"created_at,omitempty"`
}

type UserPayload struct {
	ErrorPayload
	User *User `json:"user,omitempty"`
}

func (s *UsersService) AllUpdatedSince(updatedSince time.Time) ([]*User, error) {
	peopleUrl := "/people"
	if !updatedSince.IsZero() {
		values := make(url.Values)
		values.Add("updated_since", updatedSince.UTC().String())
		peopleUrl = peopleUrl + "?" + values.Encode()
	}
	response, err := s.h.ProcessRequest("GET", peopleUrl, nil)
	if err != nil {
		return nil, err
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	userResponses := make([]*UserPayload, 0)
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

func (s *UsersService) Create(user *User) (*User, error) {
	marshaledUser, err := json.Marshal(&UserPayload{User: user})
	if err != nil {
		return nil, err
	}
	response, err := s.h.ProcessRequest("POST", "/people", bytes.NewReader(marshaledUser))
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
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			return nil, err
		}
		return nil, &ResponseError{&apiResponse}
	}
	user.Id = userId
	return user, nil
}

func (s *UsersService) ResetPassword(user *User) error {
	marshaledUser, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = s.h.ProcessRequest("POST", fmt.Sprintf("/people/%d/reset_password", user.Id), bytes.NewBuffer(marshaledUser))
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersService) Update(user *User) (*User, error) {
	marshaledUser, err := json.Marshal(&UserPayload{User: user})
	if err != nil {
		return nil, err
	}
	response, err := s.h.ProcessRequest("PUT", fmt.Sprintf("/people/%d", user.Id), bytes.NewBuffer(marshaledUser))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			return nil, err
		}
		return nil, &ResponseError{&apiResponse}
	}
	return user, nil
}

func (s *UsersService) Delete(user *User) (bool, error) {
	response, err := s.h.ProcessRequest("DELETE", fmt.Sprintf("/people/%d", user.Id), nil)
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
	response, err := s.h.ProcessRequest("POST", fmt.Sprintf("/people/%d", user.Id), nil)
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
