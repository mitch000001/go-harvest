// +build feature

package harvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

type UserService struct {
	api *JsonApi
}

func NewUserService(api *JsonApi) *UserService {
	service := UserService{api: api}
	return &service
}

func (s *UserService) All() ([]*User, error) {
	return s.AllUpdatedSince(time.Time{})
}

func (s *UserService) AllUpdatedSince(updatedSince time.Time) ([]*User, error) {
	peopleUrl := "/people"
	if !updatedSince.IsZero() {
		values := make(url.Values)
		values.Add("updated_since", updatedSince.UTC().String())
		peopleUrl = peopleUrl + "?" + values.Encode()
	}
	response, err := s.api.Process("GET", peopleUrl, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
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

func (s *UserService) Find(id int) (*User, error) {
	response, err := s.api.Process("GET", fmt.Sprintf("/people/%d", id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == 404 {
		return nil, &ResponseError{&ErrorPayload{fmt.Sprintf("No user found with id %d", id)}}
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	userPayload := UserPayload{}
	err = json.Unmarshal(responseBytes, &userPayload)
	if err != nil {
		return nil, err
	}
	return userPayload.User, nil
}

func (s *UserService) Create(user *User) (*User, error) {
	marshaledUser, err := json.Marshal(&UserPayload{User: user})
	if err != nil {
		return nil, err
	}
	response, err := s.api.Process("POST", "/people", bytes.NewReader(marshaledUser))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	userId := -1
	fmt.Printf("Headers: %+#v\n", response.Header)
	if response.StatusCode == 201 {
		location := response.Header.Get("Location")
		fmt.Sscanf(location, "/people/%d", &userId)
	}
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
	user.SetId(userId)
	return user, nil
}

func (s *UserService) Update(user *User) (*User, error) {
	marshaledUser, err := json.Marshal(&UserPayload{User: user})
	if err != nil {
		return nil, err
	}
	response, err := s.api.Process("PUT", fmt.Sprintf("/people/%d", user.Id), bytes.NewBuffer(marshaledUser))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
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

func (s *UserService) Delete(user *User) (bool, error) {
	response, err := s.api.Process("DELETE", fmt.Sprintf("/people/%d", user.Id), nil)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return true, nil
	} else if response.StatusCode == 400 {
		return false, nil
	} else {
		panic(response.Status)
	}
}

func (s *UserService) Toggle(user *User) (bool, error) {
	response, err := s.api.Process("POST", fmt.Sprintf("/people/%d", user.Id), nil)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return true, nil
	} else if response.StatusCode == 400 {
		return false, nil
	} else {
		panic(response.Status)
	}
}
