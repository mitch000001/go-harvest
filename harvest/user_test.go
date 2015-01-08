package harvest

import (
	"net/url"
	"reflect"
	"testing"
)

func TestUserServiceCalls(t *testing.T) {
	type testFunc func(*apiWrapperTestData, *bool) Api

	expectedParams := url.Values{"foo": []string{"bar"}}

	var tests = []struct {
		testData *apiWrapperTestData
		testFn   testFunc
		apiFn    string
		args     []interface{}
	}{
		{
			&apiWrapperTestData{
				expectedParams:       expectedParams,
				expectedDataType:     reflect.TypeOf(&[]*User{}),
				expectedErrorMessage: "ERR",
			},
			testApiAllWrapper,
			"All",
			[]interface{}{&[]*User{}, expectedParams},
		},
		{
			&apiWrapperTestData{
				expectedParams:       expectedParams,
				expectedIdType:       reflect.TypeOf(12),
				expectedDataType:     reflect.TypeOf(&User{}),
				expectedErrorMessage: "ERR",
			},
			testApiFindWrapper,
			"Find",
			[]interface{}{12, &User{}, expectedParams},
		},
		{
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&User{}),
				expectedErrorMessage: "ERR",
			},
			testApiCreateWrapper,
			"Create",
			[]interface{}{&User{}},
		},
		{
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&User{}),
				expectedErrorMessage: "ERR",
			},
			testApiUpdateWrapper,
			"Update",
			[]interface{}{&User{}},
		},
		{
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&User{}),
				expectedErrorMessage: "ERR",
			},
			testApiDeleteWrapper,
			"Delete",
			[]interface{}{&User{}},
		},
	}
	for _, test := range tests {
		called := false
		api := test.testFn(test.testData, &called)
		service := NewUserService(api)
		serviceValue := reflect.ValueOf(service)
		testFn := serviceValue.MethodByName(test.apiFn)
		var args []reflect.Value
		for _, v := range test.args {
			args = append(args, reflect.ValueOf(v))
		}
		res := testFn.Call(args)

		if !called {
			t.Logf("Expected Api.%s method to have been called, was not.\n", test.apiFn)
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
}

func TestUserServiceAll(t *testing.T) {
	expectedParams := url.Values{"foo": []string{"bar"}}
	testData := apiWrapperTestData{
		expectedParams:       expectedParams,
		expectedDataType:     reflect.TypeOf(&[]*User{}),
		expectedErrorMessage: "ERR",
	}
	called := false
	api := testApiAllWrapper(&testData, &called)
	service := NewUserService(api)

	var users []*User

	err := service.All(&users, expectedParams)

	if !called {
		t.Logf("Expected Api.Find method to get called\n")
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
	}

	if err != nil {
		expectedMessage := "ERR"
		actualMessage := err.Error()
		if expectedMessage != actualMessage {
			t.Logf("Expected error to have message '%q', got '%q'\n", expectedMessage, actualMessage)
			t.Fail()
		}
	}
}

func TestUserServiceFind(t *testing.T) {
	expectedParams := url.Values{"foo": []string{"bar"}}
	testData := apiWrapperTestData{
		expectedParams:       expectedParams,
		expectedIdType:       reflect.TypeOf(12),
		expectedDataType:     reflect.TypeOf(&User{}),
		expectedErrorMessage: "ERR",
	}
	called := false
	api := testApiFindWrapper(&testData, &called)
	service := NewUserService(api)

	var user User

	err := service.Find(12, &user, expectedParams)

	if !called {
		t.Logf("Expected Api.Find method to get called\n")
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
	}

	if err != nil {
		expectedMessage := "ERR"
		actualMessage := err.Error()
		if expectedMessage != actualMessage {
			t.Logf("Expected error to have message '%q', got '%q'\n", expectedMessage, actualMessage)
			t.Fail()
		}
	}
}
