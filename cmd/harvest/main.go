package main

import (
	"code.google.com/p/goauth2/oauth"
	"github.com/mitch000001/go-harvest/harvest"
)

func main() {

}

// NewBasicAuthClient creates a new Client with BasicAuth as authentication method
func NewBasicAuthClient(subdomain string, config *BasicAuthConfig) (*harvest.Harvest, error) {
	h, err := harvest.NewHarvest(subdomain)
	if err != nil {
		return nil, err
	}
	h.authenticationTransport = &Transport{Config: config}
	return h, nil
}

// NewOAuthClient creates a new Client with OAuth as authentication method
func NewOAuthClient(subdomain string, config *oauth.Config) (*harvest.Harvest, error) {
	h, err := harvest.NewHarvest(subdomain)
	if err != nil {
		return nil, err
	}
	h.authenticationTransport = &oauth.Transport{Config: config}
	return h, err
}
