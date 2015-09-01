// +build !feature

package harvest

import (
	"net/url"
	"reflect"
	"testing"
)

var (
	expectedTaskAssignmentServiceParams	= url.Values{"foo": []string{"bar"}}

	testsTaskAssignmentService	= map[string]struct {	// apiFn to testData
		testData	*apiWrapperTestData
		testFn		testFunc
		args		[]interface{}
	}{
		"All": {
			&apiWrapperTestData{
				expectedParams:		expectedTaskAssignmentServiceParams,
				expectedDataType:	reflect.TypeOf(&[]*TaskAssignment{}),
				expectedErrorMessage:	"ERR",
			},
			testApiAllWrapper,
			[]interface{}{&[]*TaskAssignment{}, expectedTaskAssignmentServiceParams},
		},
		"Find": {
			&apiWrapperTestData{
				expectedParams:		expectedTaskAssignmentServiceParams,
				expectedIdType:		reflect.TypeOf(12),
				expectedDataType:	reflect.TypeOf(&TaskAssignment{}),
				expectedErrorMessage:	"ERR",
			},
			testApiFindWrapper,
			[]interface{}{12, &TaskAssignment{}, expectedTaskAssignmentServiceParams},
		},
		"Create": {
			&apiWrapperTestData{
				expectedDataType:	reflect.TypeOf(&TaskAssignment{}),
				expectedErrorMessage:	"ERR",
			},
			testApiCreateWrapper,
			[]interface{}{&TaskAssignment{}},
		},
		"Update": {
			&apiWrapperTestData{
				expectedDataType:	reflect.TypeOf(&TaskAssignment{}),
				expectedErrorMessage:	"ERR",
			},
			testApiUpdateWrapper,
			[]interface{}{&TaskAssignment{}},
		},
		"Delete": {
			&apiWrapperTestData{
				expectedDataType:	reflect.TypeOf(&TaskAssignment{}),
				expectedErrorMessage:	"ERR",
			},
			testApiDeleteWrapper,
			[]interface{}{&TaskAssignment{}},
		},
	}
)

func TestTaskAssignmentServiceAll(t *testing.T) {
	testTaskAssignmentServiceMethod(t, "All")
}

func TestTaskAssignmentServiceFind(t *testing.T) {
	testTaskAssignmentServiceMethod(t, "Find")
}

func TestTaskAssignmentServiceCreate(t *testing.T) {
	testTaskAssignmentServiceMethod(t, "Create")
}

func TestTaskAssignmentServiceUpdate(t *testing.T) {
	testTaskAssignmentServiceMethod(t, "Update")
}

func TestTaskAssignmentServiceDelete(t *testing.T) {
	testTaskAssignmentServiceMethod(t, "Delete")
}

func testTaskAssignmentServiceMethod(t *testing.T, name string) {
	called := false
	test, ok := testsTaskAssignmentService[name]
	if !ok {
		t.Logf("No test data for method '%s' defined.\n", name)
		t.FailNow()
	}
	api := test.testFn(test.testData, &called)
	service := &TaskAssignmentService{endpoint: api}
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
