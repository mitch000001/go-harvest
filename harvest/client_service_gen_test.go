// DO NOT EDIT!
// This file is generated by the api generator.

// +build !feature

package harvest

import (
	"net/url"
	"reflect"
	"testing"
)

var (
	expectedClientServiceParams	= url.Values{"foo": []string{"bar"}}

	testsClientService	= map[string]struct {	// apiFn to testData
		testData	*apiWrapperTestData
		testFn		testFunc
		args		[]interface{}
	}{
		"All": {
			&apiWrapperTestData{
				expectedParams:		expectedClientServiceParams,
				expectedDataType:	reflect.TypeOf(&[]*Client{}),
				expectedErrorMessage:	"ERR",
			},
			testApiAllWrapper,
			[]interface{}{&[]*Client{}, expectedClientServiceParams},
		},
		"Find": {
			&apiWrapperTestData{
				expectedParams:		expectedClientServiceParams,
				expectedIdType:		reflect.TypeOf(12),
				expectedDataType:	reflect.TypeOf(&Client{}),
				expectedErrorMessage:	"ERR",
			},
			testApiFindWrapper,
			[]interface{}{12, &Client{}, expectedClientServiceParams},
		},
		"Create": {
			&apiWrapperTestData{
				expectedDataType:	reflect.TypeOf(&Client{}),
				expectedErrorMessage:	"ERR",
			},
			testApiCreateWrapper,
			[]interface{}{&Client{}},
		},
		"Update": {
			&apiWrapperTestData{
				expectedDataType:	reflect.TypeOf(&Client{}),
				expectedErrorMessage:	"ERR",
			},
			testApiUpdateWrapper,
			[]interface{}{&Client{}},
		},
		"Delete": {
			&apiWrapperTestData{
				expectedDataType:	reflect.TypeOf(&Client{}),
				expectedErrorMessage:	"ERR",
			},
			testApiDeleteWrapper,
			[]interface{}{&Client{}},
		},
		"Toggle": {
			&apiWrapperTestData{
				expectedDataType:	reflect.TypeOf(&Client{}),
				expectedErrorMessage:	"ERR",
			},
			testApiToggleWrapper,
			[]interface{}{&Client{}},
		},
	}
)

func TestClientServiceAll(t *testing.T) {
	testClientServiceMethod(t, "All")
}

func TestClientServiceFind(t *testing.T) {
	testClientServiceMethod(t, "Find")
}

func TestClientServiceCreate(t *testing.T) {
	testClientServiceMethod(t, "Create")
}

func TestClientServiceUpdate(t *testing.T) {
	testClientServiceMethod(t, "Update")
}

func TestClientServiceDelete(t *testing.T) {
	testClientServiceMethod(t, "Delete")
}

func TestClientServiceToggle(t *testing.T) {
	testClientServiceMethod(t, "Toggle")
}

func testClientServiceMethod(t *testing.T, name string) {
	called := false
	test, ok := testsClientService[name]
	if !ok {
		t.Logf("No test data for method '%s' defined.\n", name)
		t.FailNow()
	}
	api := test.testFn(test.testData, &called)
	service := &ClientService{endpoint: api}
	serviceValue := reflect.ValueOf(service)
	testFn := serviceValue.MethodByName(name)
	if !testFn.IsValid() {
		t.Logf("Expected service to have method '%s', had not.\n", name)
		t.FailNow()
	}

	var args []reflect.Value
	for _, v := range test.args {
		args = append(args, reflect.ValueOf(v))
	}
	res := testFn.Call(args)

	if !called {
		t.Logf("Expected Api.%s method to have been called, was not.\n", name)
		t.Fail()
	}

	errors := test.testData.getErrors()

	if errors != "" {
		t.Logf("Found errors:\n%s", errors)
		t.Fail()
	}

	err := res[0]
	if err.IsNil() {
		t.Logf("Expected error not to be nil\n")
		t.Fail()
	}

	if !err.IsNil() {
		expectedMessage := "ERR"
		actualMessage := err.MethodByName("Error").Call([]reflect.Value{})[0].String()
		if expectedMessage != actualMessage {
			t.Logf("Expected error to have message '%q', got '%q'\n", expectedMessage, actualMessage)
			t.Fail()
		}
	}
}
