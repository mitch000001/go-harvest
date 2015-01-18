// DO NOT EDIT!
// This file is generated by the api generator.
package harvest

import (
	"net/url"
)

type UserService struct {
	api Api
}

func NewUserService(api Api) *UserService {
	service := UserService{api: api}
	return &service
}

func (s *UserService) All(users *[]*User, params url.Values) error {
	return s.api.All(users, params)
}

func (s *UserService) Find(id int, user *User, params url.Values) error {
	return s.api.Find(id, user, params)
}

func (s *UserService) Create(user *User) error {
	return s.api.Create(user)
}

func (s *UserService) Update(user *User) error {
	return s.api.Update(user)
}

func (s *UserService) Delete(user *User) error {
	return s.api.Delete(user)
}
