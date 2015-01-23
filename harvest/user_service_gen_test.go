// DO NOT EDIT!
// This file is generated by the api generator.

// +build !integration

package harvest

import (
	"net/url"
	"reflect"
	"testing"
)

var (
	expectedUserServiceParams = url.Values{"foo": []string{"bar"}}

	testsUserService = map[string]struct { // apiFn to testData
		testData *apiWrapperTestData
		testFn   testFunc
		args     []interface{}
	}{
		"All": {
			&apiWrapperTestData{
				expectedParams:       expectedUserServiceParams,
				expectedDataType:     reflect.TypeOf(&[]*User{}),
				expectedErrorMessage: "ERR",
			},
			testApiAllWrapper,
			[]interface{}{&[]*User{}, expectedUserServiceParams},
		},
		"Find": {
			&apiWrapperTestData{
				expectedParams:       expectedUserServiceParams,
				expectedIdType:       reflect.TypeOf(12),
				expectedDataType:     reflect.TypeOf(&User{}),
				expectedErrorMessage: "ERR",
			},
			testApiFindWrapper,
			[]interface{}{12, &User{}, expectedUserServiceParams},
		},
		"Create": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&User{}),
				expectedErrorMessage: "ERR",
			},
			testApiCreateWrapper,
			[]interface{}{&User{}},
		},
		"Update": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&User{}),
				expectedErrorMessage: "ERR",
			},
			testApiUpdateWrapper,
			[]interface{}{&User{}},
		},
		"Delete": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&User{}),
				expectedErrorMessage: "ERR",
			},
			testApiDeleteWrapper,
			[]interface{}{&User{}},
		},
		"Toggle": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&User{}),
				expectedErrorMessage: "ERR",
			},
			testApiToggleWrapper,
			[]interface{}{&User{}},
		},
	}
)

func TestUserServiceAll(t *testing.T) {
	testUserServiceMethod(t, "All")
}

func TestUserServiceFind(t *testing.T) {
	testUserServiceMethod(t, "Find")
}

func TestUserServiceCreate(t *testing.T) {
	testUserServiceMethod(t, "Create")
}

func TestUserServiceUpdate(t *testing.T) {
	testUserServiceMethod(t, "Update")
}

func TestUserServiceDelete(t *testing.T) {
	testUserServiceMethod(t, "Delete")
}

func TestUserServiceToggle(t *testing.T) {
	testUserServiceMethod(t, "Toggle")
}

func testUserServiceMethod(t *testing.T, name string) {
	called := false
	test, ok := testsUserService[name]
	if !ok {
		t.Logf("No test data for method '%s' defined.\n", name)
		t.FailNow()
	}
	api := test.testFn(test.testData, &called)
	service := &UserService{endpoint: api}
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
