package auth

import (
	"net/http"
	"reflect"
	"testing"
)

func TestTransportClient(t *testing.T) {
	config := BasicAuthConfig{Username: "foo", Password: "bar"}
	transport := Transport{Config: &config}

	client := transport.Client()

	if client == nil {
		t.Logf("Expected client not to be nil\n")
		t.Fail()
	}

	if !reflect.DeepEqual(&transport, client.Transport) {
		t.Logf("Expected client transport to equal '%+#v', got '%+#v'", transport, client.Transport)
		t.Fail()
	}
}

func TestTransportRoundTrip(t *testing.T) {
	config := BasicAuthConfig{Username: "foo", Password: "bar"}
	transport := Transport{Config: &config}

	request, _ := http.NewRequest("GET", "/foo", nil)

	transport.RoundTrip(request)
}
