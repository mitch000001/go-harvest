package mock

import (
	"fmt"
	"net/url"

	"github.com/mitch000001/go-harvest/harvest"
)

type DayEntryService struct {
	Entries []*harvest.DayEntry
	harvest.CrudEndpoint
}

func (d DayEntryService) All(data interface{}, params url.Values) error {
	timeframe, err := harvest.TimeframeFromQuery(params)
	if err != nil {
		return fmt.Errorf("Error while parsing timeframe: %v", err)
	}
	entries := make([]*harvest.DayEntry, 0)
	for _, entry := range d.Entries {
		if timeframe.IsInTimeframe(entry.SpentAt) {
			entries = append(entries, entry)
		}
	}
	*(data.(*[]*harvest.DayEntry)) = entries
	return nil
}

func (d DayEntryService) Path() string {
	return "entries"
}

func (d DayEntryService) URL() url.URL {
	return url.URL{}
}

func NewDayEntryService(service DayEntryService) *harvest.DayEntryService {
	return harvest.NewDayEntryService(service)
}
