// +build feature

package harvest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ClientService struct {
	api *JsonApi
}

func NewClientService(api *JsonApi) *ClientService {
	return &ClientService{api}
}

func (c *ClientService) All() ([]*Client, error) {
	response, err := c.api.Process("GET", "/clients", nil)
	if err != nil {
		return nil, err
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	clientPayloads := make([]*ClientPayload, 0)
	err = json.Unmarshal(responseBytes, &clientPayloads)
	if err != nil {
		return nil, err
	}
	clients := make([]*Client, len(clientPayloads))
	for i, c := range clientPayloads {
		clients[i] = c.Client
	}
	return clients, nil
}

func (c *ClientService) Find(id int) (*Client, error) {
	response, err := c.api.Process("GET", fmt.Sprintf("/clients/%d", id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == 404 {
		return nil, &ResponseError{&ErrorPayload{fmt.Sprintf("No client found for id %d", id)}}
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	clientPayload := ClientPayload{}
	err = json.Unmarshal(responseBytes, &clientPayload)
	if err != nil {
		return nil, err
	}
	return clientPayload.Client, nil
}
