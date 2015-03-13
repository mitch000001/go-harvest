package auth

import (
	"net/http"
)

type BasicAuthError struct {
	prefix string
	msg    string
}

func (bae BasicAuthError) Error() string {
	return "BasicAuthError: " + bae.prefix + ": " + bae.msg
}

// BasicAuthConfig is the configuration of a basic auth consumer
type BasicAuthConfig struct {
	Username string
	Password string
}

// Transport implements http.RoundTripper. When configured with a valid
// BasicAuthConfig it can be used to make authenticated HTTP requests.
//
//	t := &Transport{config}
//	r, err := t.Client().Get("http://example.org/url/requiring/auth")
//
// It will automatically refresh the Token if it can,
// updating the supplied Token in place.
type Transport struct {
	Config *BasicAuthConfig

	// Transport is the HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	// (It should never be an oauth.Transport.)
	Transport http.RoundTripper
}

// Client returns an *http.Client that makes OAuth-authenticated requests.
func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}

// Fetches the internal transport.
func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// RoundTrip executes a single HTTP transaction using the Transport's
// Token as authorization headers.
//
// This method will attempt to renew the Token if it has expired and may return
// an error related to that Token renewal before attempting the client request.
// If the Token cannot be renewed a non-nil os.Error value will be returned.
// If the Token is invalid callers should expect HTTP-level errors,
// as indicated by the Response's StatusCode.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Config == nil {
		return nil, BasicAuthError{"RoundTrip", "no Config supplied"}
	}
	// To set the Basic Auth, we must make a copy of the Request
	// so that we don't modify the Request we were given.
	// This is required by the specification of http.RoundTripper.
	req = cloneRequest(req)
	req.SetBasicAuth(t.Config.Username, t.Config.Password)

	// Make the HTTP request.
	return t.transport().RoundTrip(req)
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}
