package auth

import (
	"net/http"

	"code.google.com/p/goauth2/oauth"
	"github.com/mitch000001/go-harvest/harvest"
)

// NewBasicAuthClient creates a new ClientProvider with BasicAuth as authentication method
func NewBasicAuthClientProvider(config *BasicAuthConfig) harvest.HttpClientProvider {
	clientProvider := &Transport{Config: config}
	return ClientProviderFunc(clientProvider.Client)
}

// NewOAuthClient creates a new ClientProvider with OAuth as authentication method
func NewOAuthClientProvider(config *oauth.Config) harvest.HttpClientProvider {
	clientProvider := &oauth.Transport{Config: config}
	return ClientProviderFunc(clientProvider.Client)
}

type ClientProviderFunc func() *http.Client

func (cf ClientProviderFunc) Client() harvest.HttpClient {
	return cf()
}
