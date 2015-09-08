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
			&harvest.DayEntry{ID: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 1, time.UTC)},
			&harvest.DayEntry{ID: 2, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 2, 1, time.UTC)},
			&harvest.DayEntry{ID: 3, Hours: 8, TaskId: 5, SpentAt: harvest.Date(2015, 1, 19, time.UTC)},
			&harvest.DayEntry{ID: 4, Hours: 8, TaskId: 7, SpentAt: harvest.Date(2015, 1, 20, time.UTC)},
			&harvest.DayEntry{ID: 5, Hours: 8, TaskId: 9, SpentAt: harvest.Date(2015, 1, 21, time.UTC)},
		},
		BillableTasks: []int{2, 5},
	}

	var entries []*harvest.DayEntry
	var params harvest.Params
	timeframe := harvest.NewTimeframe(2015, 1, 1, 2015, 4, 1, time.UTC)

	err := service.All(&entries, params.ForTimeframe(timeframe).Values())

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedEntries := []*harvest.DayEntry{
		&harvest.DayEntry{ID: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 1, time.UTC)},
		&harvest.DayEntry{ID: 2, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 2, 1, time.UTC)},
		&harvest.DayEntry{ID: 3, Hours: 8, TaskId: 5, SpentAt: harvest.Date(2015, 1, 19, time.UTC)},
		&harvest.DayEntry{ID: 4, Hours: 8, TaskId: 7, SpentAt: harvest.Date(2015, 1, 20, time.UTC)},
		&harvest.DayEntry{ID: 5, Hours: 8, TaskId: 9, SpentAt: harvest.Date(2015, 1, 21, time.UTC)},
	}

	if !reflect.DeepEqual(expectedEntries, entries) {
		t.Logf("Expected entries to equal\n%q\n\tgot\n%q\n", expectedEntries, entries)
		t.Fail()
	}

	// Proper filtering for timeframes
	timeframe = harvest.NewTimeframe(2015, 1, 1, 2015, 1, 25, time.UTC)
	params = harvest.Params{}

	err = service.All(&entries, params.ForTimeframe(timeframe).Values())

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedEntries = []*harvest.DayEntry{
		&harvest.DayEntry{ID: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 1, time.UTC)},
		&harvest.DayEntry{ID: 3, Hours: 8, TaskId: 5, SpentAt: harvest.Date(2015, 1, 19, time.UTC)},
		&harvest.DayEntry{ID: 4, Hours: 8, TaskId: 7, SpentAt: harvest.Date(2015, 1, 20, time.UTC)},
		&harvest.DayEntry{ID: 5, Hours: 8, TaskId: 9, SpentAt: harvest.Date(2015, 1, 21, time.UTC)},
	}

	if !reflect.DeepEqual(expectedEntries, entries) {
		t.Logf("Expected entries to equal\n%q\n\tgot\n%q\n", expectedEntries, entries)
		t.Fail()
	}

	// proper filtering for billable
	params = harvest.Params{}

	err = service.All(&entries, params.ForTimeframe(timeframe).Billable(true).Values())

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedEntries = []*harvest.DayEntry{
		&harvest.DayEntry{ID: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 1, time.UTC)},
		&harvest.DayEntry{ID: 3, Hours: 8, TaskId: 5, SpentAt: harvest.Date(2015, 1, 19, time.UTC)},
	}

	if !reflect.DeepEqual(expectedEntries, entries) {
		t.Logf("Expected entries to equal\n%q\n\tgot\n%q\n", expectedEntries, entries)
		t.Fail()
	}

	// proper filtering for nonbillable
	params = harvest.Params{}

	err = service.All(&entries, params.ForTimeframe(timeframe).Billable(false).Values())

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedEntries = []*harvest.DayEntry{
		&harvest.DayEntry{ID: 4, Hours: 8, TaskId: 7, SpentAt: harvest.Date(2015, 1, 20, time.UTC)},
		&harvest.DayEntry{ID: 5, Hours: 8, TaskId: 9, SpentAt: harvest.Date(2015, 1, 21, time.UTC)},
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

func TestDayEntryFilterAdd(t *testing.T) {
	filter := dayEntryFilter{}

	filterOdd := func(e *harvest.DayEntry) bool {
		return e.ID%2 != 0
	}

	filter.add(filterOdd)

	if len(filter) != 1 {
		t.Logf("Expected filter to have one item, got %d\n", len(filter))
		t.Fail()
	}
}

func TestDayEntryFilterApply(t *testing.T) {
	filter := dayEntryFilter{}

	filterOdd := func(e *harvest.DayEntry) bool {
		return e.ID%2 != 0
	}
	filterMod3 := func(e *harvest.DayEntry) bool {
		return e.ID%3 != 0
	}

	filter.add(filterOdd)
	filter.add(filterMod3)

	dataSet := []*harvest.DayEntry{
		&harvest.DayEntry{ID: 1},
		&harvest.DayEntry{ID: 2},
		&harvest.DayEntry{ID: 3},
		&harvest.DayEntry{ID: 4},
		&harvest.DayEntry{ID: 5},
	}

	expectedDataSet := []*harvest.DayEntry{
		&harvest.DayEntry{ID: 1},
		&harvest.DayEntry{ID: 5},
	}

	var result []*harvest.DayEntry

	for _, entry := range dataSet {
		if filter.apply(entry) {
			result = append(result, entry)
		}
	}

	if !reflect.DeepEqual(expectedDataSet, result) {
		t.Logf("Expected result to equal\n%q\n\tgot\n%q\n", expectedDataSet, result)
		t.Fail()
	}
}
