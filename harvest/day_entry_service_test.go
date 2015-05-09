package harvest

import (
	"net/url"
	"reflect"
	"testing"
)

func TestNewDayEntryService(t *testing.T) {

}

func TestDayEntryServiceAll(t *testing.T) {
	// Test proper delegation
	called := false
	expectedServiceParams := Params{}
	expectedServiceParams.Add("from", "foo")
	expectedServiceParams.Add("to", "bar")
	testData := apiWrapperTestData{
		expectedDataType:     reflect.TypeOf(&[]*DayEntry{}),
		expectedParams:       expectedServiceParams.Values(),
		expectedErrorMessage: "FooBar",
	}

	testAllWrapper := testApiAllWrapper(&testData, &called)

	service := NewDayEntryService(testAllWrapper)

	var dayEntries []*DayEntry

	err := service.All(&dayEntries, expectedServiceParams.Values())

	if !called {
		t.Logf("Expected Api.All method to have been called, was not.\n")
		t.Fail()
	}

	errors := testData.getErrors()

	if errors != "" {
		t.Logf("Found errors:\n%s", errors)
		t.Fail()
	}

	if err == nil {
		t.Logf("Expected error not to be nil\n")
		t.Fail()
	} else {
		expectedMessage := "FooBar"
		actualMessage := err.Error()
		if expectedMessage != actualMessage {
			t.Logf("Expected error to have message '%q', got '%q'\n", expectedMessage, actualMessage)
			t.Fail()
		}
	}

	// Test missing params
	endpoint := testApiAll(func(interface{}, url.Values) error { return nil })

	service = NewDayEntryService(endpoint)

	var tests = []struct {
		params       Params
		expectError  bool
		errorMessage string
	}{
		{
			Params{},
			true,
			"Bad Request: 'from' and 'to' query parameter are not optional!",
		}, {
			Params{"from": []string{"foo"}},
			true,
			"Bad Request: 'from' and 'to' query parameter are not optional!",
		}, {
			Params{"to": []string{"bar"}},
			true,
			"Bad Request: 'from' and 'to' query parameter are not optional!",
		}, {
			Params{"from": []string{"foo"}, "to": []string{"bar"}},
			false,
			"",
		},
	}
	for _, test := range tests {
		err := service.All(&[]*DayEntry{}, test.params.Values())
		if test.expectError {
			if err == nil {
				t.Logf("Expected error, got nil\n")
				t.Fail()
			} else {
				msg := err.Error()
				if msg != test.errorMessage {
					t.Logf("Expected error message %q, got %q\n", test.errorMessage, msg)
					t.Fail()
				}
			}
		} else {
			if err != nil {
				t.Logf("Expected no error, got %T: %v\n", err, err)
				t.Fail()
			}
		}
	}
}
