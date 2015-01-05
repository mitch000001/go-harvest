package main

import (
	"net/http"

	"code.google.com/p/goauth2/oauth"
	"github.com/mitch000001/go-harvest/harvest"
)

func main() {

}

// NewBasicAuthClient creates a new Client with BasicAuth as authentication method
func NewBasicAuthClient(subdomain string, config *BasicAuthConfig) (*harvest.Harvest, error) {
	clientProvider := &Transport{Config: config}
	h, err := harvest.NewHarvest(subdomain, BuildClientProvider(clientProvider.Client))
	if err != nil {
		return nil, err
	}
	return h, nil
}

// NewOAuthClient creates a new Client with OAuth as authentication method
func NewOAuthClient(subdomain string, config *oauth.Config) (*harvest.Harvest, error) {
	clientProvider := &oauth.Transport{Config: config}
	h, err := harvest.NewHarvest(subdomain, BuildClientProvider(clientProvider.Client))
	if err != nil {
		return nil, err
	}
	return h, err
}

type clientProviderFunc func() *http.Client

func (cf clientProviderFunc) Client() harvest.HttpClient {
	return cf()
}

type clientProviderWrapper struct {
	clientProviderFunc clientProviderFunc
}

func (cw *clientProviderWrapper) Client() harvest.HttpClient {
	return cw.clientProviderFunc()
}

func BuildClientProvider(f clientProviderFunc) harvest.HttpClientProvider {
	wrapper := &clientProviderWrapper{clientProviderFunc: f}
	return wrapper
}
