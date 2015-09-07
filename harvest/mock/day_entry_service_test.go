package mock

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/mitch000001/go-harvest/harvest"
)

func TestNewDayEntryService(t *testing.T) {
	service := DayEntryService{
		Entries: []*harvest.DayEntry{
			&harvest.DayEntry{ID: 1, Hours: 8, TaskId: 2},
		},
	}

	dayEntryService := NewDayEntryService(service)

	if dayEntryService == nil {
		t.Logf("Expected service not to be nil\n")
		t.Fail()
	}

	var entries []*harvest.DayEntry
	timeframe := harvest.NewTimeframe(2015, 1, 1, 2015, 4, 1, time.UTC)
	var params harvest.Params

	err := dayEntryService.All(&entries, params.ForTimeframe(timeframe).Values())

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}
}

func TestDayEntryServiceAll(t *testing.T) {
	service := DayEntryService{
		Entries: []*harvest.DayEntry{
			&harvest.DayEntry{ID: 1, Hours: 8, TaskId: 2},
		},
	}

	var entries []*harvest.DayEntry

	err := service.All(&entries, nil)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedEntries := []*harvest.DayEntry{
		&harvest.DayEntry{ID: 1, Hours: 8, TaskId: 2},
	}

	if !reflect.DeepEqual(expectedEntries, entries) {
		t.Logf("Expected entries to equal\n%q\n\tgot\n%q\n", expectedEntries, entries)
		t.Fail()
	}
}

func TestDayEntryServicePath(t *testing.T) {
	service := DayEntryService{}

	path := service.Path()

	if path != "entries" {
		t.Logf("Expected Path to return 'entries', got %q\n", path)
		t.Fail()
	}
}

func TestDayEntryServiceURL(t *testing.T) {
	service := DayEntryService{}

	actualUrl := service.URL()

	expectedUrl := url.URL{}

	if expectedUrl.String() != actualUrl.String() {
		t.Logf("Expected URL to return\n%q\n\tgot\n%q\n", expectedUrl, actualUrl)
		t.Fail()
	}
}
