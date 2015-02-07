package harvest

import (
	"fmt"
	"net/url"
)

type DayEntryService struct {
	endpoint AllEndpoint
}

func NewDayEntryService(endpoint AllEndpoint) *DayEntryService {
	return &DayEntryService{endpoint: endpoint}
}

func (d *DayEntryService) All(dayEntries *[]*DayEntry, params url.Values) error {
	if len(params) == 0 || params.Get("from") == "" || params.Get("to") == "" {
		return fmt.Errorf("Bad Request: 'from' and 'to' query parameter are not optional!")
	}
	return d.endpoint.All(dayEntries, params)
}
