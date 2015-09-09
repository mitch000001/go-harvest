package mock

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/mitch000001/go-harvest/harvest"
)

func TestNewDayEntryService(t *testing.T) {
	endpoint := DayEntryEndpoint{
		Entries: []*harvest.DayEntry{
			&harvest.DayEntry{ID: 1, UserId: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 2, time.UTC)},
			&harvest.DayEntry{ID: 1, UserId: 2, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 3, time.UTC)},
		},
		UserId: 1,
	}

	dayEntryService := NewDayEntryService(&endpoint)

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

	expectedEntries := []*harvest.DayEntry{
		&harvest.DayEntry{ID: 1, UserId: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 2, time.UTC)},
	}

	if !reflect.DeepEqual(expectedEntries, entries) {
		t.Logf("Expected entries to equal\n%q\n\tgot\n%q\n", expectedEntries, entries)
		t.Fail()
	}
}

func TestDayEntryEndpointAll(t *testing.T) {
	endpoint := DayEntryEndpoint{
		Entries: []*harvest.DayEntry{
			&harvest.DayEntry{ID: 1, UserId: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 1, time.UTC)},
			&harvest.DayEntry{ID: 2, UserId: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 2, 1, time.UTC)},
			&harvest.DayEntry{ID: 3, UserId: 1, Hours: 8, TaskId: 5, SpentAt: harvest.Date(2015, 1, 19, time.UTC)},
			&harvest.DayEntry{ID: 4, UserId: 1, Hours: 8, TaskId: 7, SpentAt: harvest.Date(2015, 1, 20, time.UTC)},
			&harvest.DayEntry{ID: 5, UserId: 1, Hours: 8, TaskId: 9, SpentAt: harvest.Date(2015, 1, 21, time.UTC)},
			&harvest.DayEntry{ID: 11, UserId: 2, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 1, time.UTC)},
			&harvest.DayEntry{ID: 12, UserId: 2, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 2, 1, time.UTC)},
			&harvest.DayEntry{ID: 13, UserId: 2, Hours: 8, TaskId: 5, SpentAt: harvest.Date(2015, 1, 19, time.UTC)},
			&harvest.DayEntry{ID: 14, UserId: 2, Hours: 8, TaskId: 7, SpentAt: harvest.Date(2015, 1, 20, time.UTC)},
			&harvest.DayEntry{ID: 15, UserId: 2, Hours: 8, TaskId: 9, SpentAt: harvest.Date(2015, 1, 21, time.UTC)},
		},
		BillableTasks: []int{2, 5},
		UserId:        1,
	}

	var entries []*harvest.DayEntry
	var params harvest.Params
	timeframe := harvest.NewTimeframe(2015, 1, 1, 2015, 4, 1, time.UTC)

	err := endpoint.All(&entries, params.ForTimeframe(timeframe).Values())

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedEntries := []*harvest.DayEntry{
		&harvest.DayEntry{ID: 1, UserId: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 1, time.UTC)},
		&harvest.DayEntry{ID: 2, UserId: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 2, 1, time.UTC)},
		&harvest.DayEntry{ID: 3, UserId: 1, Hours: 8, TaskId: 5, SpentAt: harvest.Date(2015, 1, 19, time.UTC)},
		&harvest.DayEntry{ID: 4, UserId: 1, Hours: 8, TaskId: 7, SpentAt: harvest.Date(2015, 1, 20, time.UTC)},
		&harvest.DayEntry{ID: 5, UserId: 1, Hours: 8, TaskId: 9, SpentAt: harvest.Date(2015, 1, 21, time.UTC)},
	}

	if !reflect.DeepEqual(expectedEntries, entries) {
		t.Logf("Expected entries to equal\n%q\n\tgot\n%q\n", expectedEntries, entries)
		t.Fail()
	}

	// Proper filtering for timeframes
	timeframe = harvest.NewTimeframe(2015, 1, 1, 2015, 1, 25, time.UTC)
	params = harvest.Params{}

	err = endpoint.All(&entries, params.ForTimeframe(timeframe).Values())

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedEntries = []*harvest.DayEntry{
		&harvest.DayEntry{ID: 1, UserId: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 1, time.UTC)},
		&harvest.DayEntry{ID: 3, UserId: 1, Hours: 8, TaskId: 5, SpentAt: harvest.Date(2015, 1, 19, time.UTC)},
		&harvest.DayEntry{ID: 4, UserId: 1, Hours: 8, TaskId: 7, SpentAt: harvest.Date(2015, 1, 20, time.UTC)},
		&harvest.DayEntry{ID: 5, UserId: 1, Hours: 8, TaskId: 9, SpentAt: harvest.Date(2015, 1, 21, time.UTC)},
	}

	if !reflect.DeepEqual(expectedEntries, entries) {
		t.Logf("Expected entries to equal\n%q\n\tgot\n%q\n", expectedEntries, entries)
		t.Fail()
	}

	// proper filtering for billable
	params = harvest.Params{}

	err = endpoint.All(&entries, params.ForTimeframe(timeframe).Billable(true).Values())

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedEntries = []*harvest.DayEntry{
		&harvest.DayEntry{ID: 1, UserId: 1, Hours: 8, TaskId: 2, SpentAt: harvest.Date(2015, 1, 1, time.UTC)},
		&harvest.DayEntry{ID: 3, UserId: 1, Hours: 8, TaskId: 5, SpentAt: harvest.Date(2015, 1, 19, time.UTC)},
	}

	if !reflect.DeepEqual(expectedEntries, entries) {
		t.Logf("Expected entries to equal\n%q\n\tgot\n%q\n", expectedEntries, entries)
		t.Fail()
	}

	// proper filtering for nonbillable
	params = harvest.Params{}

	err = endpoint.All(&entries, params.ForTimeframe(timeframe).Billable(false).Values())

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedEntries = []*harvest.DayEntry{
		&harvest.DayEntry{ID: 4, UserId: 1, Hours: 8, TaskId: 7, SpentAt: harvest.Date(2015, 1, 20, time.UTC)},
		&harvest.DayEntry{ID: 5, UserId: 1, Hours: 8, TaskId: 9, SpentAt: harvest.Date(2015, 1, 21, time.UTC)},
	}

	if !reflect.DeepEqual(expectedEntries, entries) {
		t.Logf("Expected entries to equal\n%q\n\tgot\n%q\n", expectedEntries, entries)
		t.Fail()
	}
}

func TestDayEntryEndpointPath(t *testing.T) {
	endpoint := DayEntryEndpoint{
		UserId: 1,
	}

	path := endpoint.Path()

	if path != "/1/entries" {
		t.Logf("Expected Path to return '/1/entries', got %q\n", path)
		t.Fail()
	}
}

func TestDayEntryEndpointURL(t *testing.T) {
	endpoint := DayEntryEndpoint{}

	actualUrl := endpoint.URL()

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
