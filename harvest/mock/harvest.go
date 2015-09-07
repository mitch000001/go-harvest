package mock

import (
	"net/http"

	"github.com/mitch000001/go-harvest/harvest"
)

type Mock struct {
}

func New(mock Mock) (*harvest.Harvest, error) {
	client, err := harvest.New("mock", func() harvest.HttpClient { return http.DefaultClient })
	if err != nil {
		return nil, err
	}
	return client, nil
}
