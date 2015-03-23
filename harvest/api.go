package harvest

import (
	"io"
	"net/http"
	"net/url"
)

// HttpClientProvider yields a function to provide an HttpClient.
type HttpClientProvider interface {
	// Client returns an HttpClient, which defined the minimal interface
	// of a http client usable by the harvest client to process http request
	Client() HttpClient
}

// HttpClient is the minimal interface which is used by the harvest client.
type HttpClient interface {
	// Do accepts an *http.Request and processes it
	//
	// See http.Client for a possible implementation
	Do(*http.Request) (*http.Response, error)
}

type ClientProviderFunc func() *http.Client

func (cf ClientProviderFunc) Client() HttpClient {
	return cf()
}

type RequestProcessor interface {
	Process(method string, path string, body io.Reader) (*http.Response, error)
}

type CrudEndpointProvider interface {
	CrudEndpoint(string) CrudEndpoint
}

type TogglerEndpointProvider interface {
	TogglerEndpoint(string) TogglerEndpoint
}

type CrudTogglerEndpointProvider interface {
	CrudTogglerEndpoint(string) CrudTogglerEndpoint
}

type Endpoint interface {
	URL() url.URL
	Path() string
}

type All interface {
	All(interface{}, url.Values) error
}

type AllEndpoint interface {
	All
	Endpoint
}

type CrudEndpoint interface {
	Crud
	Endpoint
}

type Crud interface {
	All
	Find(interface{}, interface{}, url.Values) error
	Create(CrudModel) error
	Update(CrudModel) error
	Delete(CrudModel) error
}

type TogglerEndpoint interface {
	Endpoint
	Toggler
}

type Toggler interface {
	Toggle(ActiveTogglerCrudModel) error
}

type CrudTogglerEndpoint interface {
	Endpoint
	Crud
	Toggler
}

type ActiveToggler interface {
	// Implementations of ToggleActive should toggle their active state and
	// return the current status
	ToggleActive() bool
}

type CrudModel interface {
	Type() string
	Id() int
	SetId(int)
}

type ActiveTogglerCrudModel interface {
	ActiveToggler
	CrudModel
}
