package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type HarvestClient struct {
	client *http.Client
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

func (h *HarvestClient) All(updatedSince time.Time) ([]*User, error) {
	values := url.Values{"updated-since": {updatedSince.UTC().String()}}
	response, err := h.client.Get("/people" + values.Encode())
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
