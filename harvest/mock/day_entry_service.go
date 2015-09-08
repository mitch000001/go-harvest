package mock

import (
	"fmt"
	"net/url"

	"github.com/mitch000001/go-harvest/harvest"
)

type DayEntryService struct {
	Entries       []*harvest.DayEntry
	BillableTasks []int
	harvest.CrudEndpoint
}

func (d DayEntryService) All(data interface{}, params url.Values) error {
	timeframe, err := harvest.TimeframeFromQuery(params)
	if err != nil {
		return fmt.Errorf("Error while parsing timeframe: %v", err)
	}
	filter := dayEntryFilter{}
	timeframeFilter := func(e *harvest.DayEntry) bool {
		return timeframe.IsInTimeframe(e.SpentAt)
	}
	filter.add(timeframeFilter)
	billableParam := params.Get("billable")
	if billableParam != "" {
		var billingFilter func(*harvest.DayEntry) bool
		if billableParam == "yes" {
			billingFilter = func(e *harvest.DayEntry) bool {
				res := false
				for _, taskId := range d.BillableTasks {
					if e.TaskId == taskId {
						res = true
					}
				}
				return res
			}
		} else if billableParam == "no" {
			billingFilter = func(e *harvest.DayEntry) bool {
				res := true
				for _, taskId := range d.BillableTasks {
					if e.TaskId == taskId {
						res = false
					}
				}
				return res
			}
		} else {
			return fmt.Errorf("Malformed billable param: %s", billableParam)
		}
		filter.add(billingFilter)
	}
	entries := make([]*harvest.DayEntry, 0)
	for _, entry := range d.Entries {
		if filter.apply(entry) {
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

type dayEntryFilter []func(*harvest.DayEntry) bool

func (d *dayEntryFilter) add(fn func(*harvest.DayEntry) bool) {
	*d = append(*d, fn)
}

func (d dayEntryFilter) apply(e *harvest.DayEntry) bool {
	result := true
	for _, fn := range d {
		result = result && fn(e)
	}
	return result
}
