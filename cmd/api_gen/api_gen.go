package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var (
	service          = new(serviceType)
	fileTemplate     = template.Must(template.New("file").Parse(fileTemplateContent))
	testfileTemplate = template.Must(template.New("testfile").Parse(testFileContent))
	fileNameTmpl     = "%s_service_gen.go"
	testfileNameTmpl = "%s_service_gen_test.go"
)

func init() {
	flag.Var(service, "type", `-type="Type"`)
}

func main() {
	flag.Parse()
	if service == nil {
		fmt.Printf("No service type given. Aborting...\n")
		os.Exit(1)
	}
	fname := fmt.Sprintf(fileNameTmpl, service.Param)
	file, err := os.Create(fname)
	if err != nil {
		fmt.Printf("There was an error creating the file: %s\n", fname)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	err = fileTemplate.Execute(file, service)
	if err != nil {
		fmt.Printf("There was an error parsing the given file template.\n")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fname = fmt.Sprintf(testfileNameTmpl, service.Param)
	file, err = os.Create(fname)
	if err != nil {
		fmt.Printf("There was an error creating the file: %s\n", fname)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	err = testfileTemplate.Execute(file, service.Type)
	if err != nil {
		fmt.Printf("There was an error parsing the given file template.\n")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

type serviceType struct {
	Type  string
	Param string
}

func (s *serviceType) String() string {
	return s.Type
}

func (s *serviceType) Set(in string) error {
	*s = serviceType{
		Type:  in,
		Param: strings.ToLower(in),
	}
	return nil
}

var fileTemplateContent = `// DO NOT EDIT!
// This file is generated by the api generator.
package harvest

import (
	"net/url"
)

type {{.Type}}Service struct {
	endpoint CrudTogglerEndpoint
}

func New{{.Type}}Service(endpoint CrudTogglerEndpoint) *{{.Type}}Service {
	service := {{.Type}}Service{endpoint: endpoint}
	return &service
}

func (s *{{.Type}}Service) All({{.Param}}s *[]*{{.Type}}, params url.Values) error {
	return s.endpoint.All({{.Param}}s, params)
}

func (s *{{.Type}}Service) Find(id int, {{.Param}} *{{.Type}}, params url.Values) error {
	return s.endpoint.Find(id, {{.Param}}, params)
}

func (s *{{.Type}}Service) Create({{.Param}} *{{.Type}}) error {
	return s.endpoint.Create({{.Param}})
}

func (s *{{.Type}}Service) Update({{.Param}} *{{.Type}}) error {
	return s.endpoint.Update({{.Param}})
}

func (s *{{.Type}}Service) Delete({{.Param}} *{{.Type}}) error {
	return s.endpoint.Delete({{.Param}})
}

func (s *{{.Type}}Service) Toggle({{.Param}} *{{.Type}}) error {
	return s.endpoint.Toggle({{.Param}})
}
`

var testFileContent = `// DO NOT EDIT!
// This file is generated by the api generator.
package harvest

import (
	"net/url"
	"reflect"
	"testing"
)

var (
	expected{{.}}ServiceParams = url.Values{"foo": []string{"bar"}}

	tests{{.}}Service = map[string]struct { // apiFn to testData
		testData *apiWrapperTestData
		testFn   testFunc
		args     []interface{}
	}{
		"All": {
			&apiWrapperTestData{
				expectedParams:       expected{{.}}ServiceParams,
				expectedDataType:     reflect.TypeOf(&[]*{{.}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiAllWrapper,
			[]interface{}{&[]*{{.}}{}, expected{{.}}ServiceParams},
		},
		"Find": {
			&apiWrapperTestData{
				expectedParams:       expected{{.}}ServiceParams,
				expectedIdType:       reflect.TypeOf(12),
				expectedDataType:     reflect.TypeOf(&{{.}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiFindWrapper,
			[]interface{}{12, &{{.}}{}, expected{{.}}ServiceParams},
		},
		"Create": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&{{.}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiCreateWrapper,

			[]interface{}{&{{.}}{}},
		},
		"Update": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&{{.}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiUpdateWrapper,

			[]interface{}{&{{.}}{}},
		},
		"Delete": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&{{.}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiDeleteWrapper,

			[]interface{}{&{{.}}{}},
		},
		"Toggle": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&{{.}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiToggleWrapper,

			[]interface{}{&{{.}}{}},
		},
	}
)

func Test{{.}}ServiceAll(t *testing.T) {
	test{{.}}ServiceMethod(t, "All")
}

func Test{{.}}ServiceFind(t *testing.T) {
	test{{.}}ServiceMethod(t, "Find")
}

func Test{{.}}ServiceCreate(t *testing.T) {
	test{{.}}ServiceMethod(t, "Create")
}

func Test{{.}}ServiceUpdate(t *testing.T) {
	test{{.}}ServiceMethod(t, "Update")
}

func Test{{.}}ServiceDelete(t *testing.T) {
	test{{.}}ServiceMethod(t, "Delete")
}

func Test{{.}}ServiceToggle(t *testing.T) {
	test{{.}}ServiceMethod(t, "Toggle")
}

func test{{.}}ServiceMethod(t *testing.T, name string) {
	called := false
	test, ok := tests{{.}}Service[name]
	if !ok {
		t.Logf("No test data for method '%s' defined.\n", name)
		t.FailNow()
	}
	api := test.testFn(test.testData, &called)
	service := New{{.}}Service(api)
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
`
