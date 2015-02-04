// +build !feature

package harvest

import (
	"net/url"
	"reflect"
	"testing"
)

var (
	expectedUserAssignmentServiceParams = url.Values{"foo": []string{"bar"}}

	testsUserAssignmentService = map[string]struct { // apiFn to testData
		testData *apiWrapperTestData
		testFn   testFunc
		args     []interface{}
	}{
		"All": {
			&apiWrapperTestData{
				expectedParams:       expectedUserAssignmentServiceParams,
				expectedDataType:     reflect.TypeOf(&[]*UserAssignment{}),
				expectedErrorMessage: "ERR",
			},
			testApiAllWrapper,
			[]interface{}{&[]*UserAssignment{}, expectedUserAssignmentServiceParams},
		},
		"Find": {
			&apiWrapperTestData{
				expectedParams:       expectedUserAssignmentServiceParams,
				expectedIdType:       reflect.TypeOf(12),
				expectedDataType:     reflect.TypeOf(&UserAssignment{}),
				expectedErrorMessage: "ERR",
			},
			testApiFindWrapper,
			[]interface{}{12, &UserAssignment{}, expectedUserAssignmentServiceParams},
		},
		"Create": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&UserAssignment{}),
				expectedErrorMessage: "ERR",
			},
			testApiCreateWrapper,
			[]interface{}{&UserAssignment{}},
		},
		"Update": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&UserAssignment{}),
				expectedErrorMessage: "ERR",
			},
			testApiUpdateWrapper,
			[]interface{}{&UserAssignment{}},
		},
		"Delete": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&UserAssignment{}),
				expectedErrorMessage: "ERR",
			},
			testApiDeleteWrapper,
			[]interface{}{&UserAssignment{}},
		},
	}
)

func TestUserAssignmentServiceAll(t *testing.T) {
	testUserAssignmentServiceMethod(t, "All")
}

func TestUserAssignmentServiceFind(t *testing.T) {
	testUserAssignmentServiceMethod(t, "Find")
}

func TestUserAssignmentServiceCreate(t *testing.T) {
	testUserAssignmentServiceMethod(t, "Create")
}

func TestUserAssignmentServiceUpdate(t *testing.T) {
	testUserAssignmentServiceMethod(t, "Update")
}

func TestUserAssignmentServiceDelete(t *testing.T) {
	testUserAssignmentServiceMethod(t, "Delete")
}

func testUserAssignmentServiceMethod(t *testing.T, name string) {
	called := false
	test, ok := testsUserAssignmentService[name]
	if !ok {
		t.Logf("No test data for method '%s' defined.\n", name)
		t.FailNow()
	}
	api := test.testFn(test.testData, &called)
	service := &UserAssignmentService{endpoint: api}
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
