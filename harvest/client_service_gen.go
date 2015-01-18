// DO NOT EDIT!
// This file is generated by the api generator.
package harvest

import (
	"net/url"
)

type ClientService struct {
	api CrudTogglerApi
}

func NewClientService(api CrudTogglerApi) *ClientService {
	service := ClientService{api: api}
	return &service
}

func (s *ClientService) All(users *[]*Client, params url.Values) error {
	return s.api.All(users, params)
}

func (s *ClientService) Find(id int, user *Client, params url.Values) error {
	return s.api.Find(id, user, params)
}

func (s *ClientService) Create(user *Client) error {
	return s.api.Create(user)
}

func (s *ClientService) Update(user *Client) error {
	return s.api.Update(user)
}

func (s *ClientService) Delete(user *Client) error {
	return s.api.Delete(user)
}

func (s *ClientService) Toggle(user *Client) error {
	return s.api.Toggle(user)
}
