package mock

import (
	"net/url"

	"github.com/mitch000001/go-harvest/harvest"
)

type DayEntryService struct {
	Entries []*harvest.DayEntry
}

func (d DayEntryService) All(entries interface{}, params url.Values) error {
	*(entries.(*[]*harvest.DayEntry)) = d.Entries
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
