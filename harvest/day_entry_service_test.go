package harvest

import (
	"reflect"
	"testing"
)

func TestNewDayEntryService(t *testing.T) {

}

func TestDayEntryServiceAll(t *testing.T) {
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
}
