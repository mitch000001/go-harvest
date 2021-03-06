package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var (
	serviceType      string
	payloadType      string
	fields           []string
	crudFlag         bool
	togglerFlag      bool
	scaffoldFlag     bool
	fileTemplate     = template.New("file")
	crudTemplate     = template.Must(fileTemplate.New("crud").Parse(crudTemplateContent))
	testfileTemplate = template.Must(template.New("testfile").Parse(testFileContent))
	fileNameTmpl     = "%s_service_gen.go"
	testfileNameTmpl = "%s_service_gen_test.go"
)

func init() {
	fileTemplate = template.Must(fileTemplate.Parse(fileTemplateContent))
	flag.BoolVar(&crudFlag, "c", true, "-c")
	flag.BoolVar(&togglerFlag, "t", false, "-t")
	flag.BoolVar(&scaffoldFlag, "s", false, "-s")
	flag.StringVar(&serviceType, "type", "", `-type="Type"`)
	flag.StringVar(&payloadType, "payload", "", `-payload="PayloadType"`)
	flag.Var((*stringsFlag)(&fields), "fields", "-fields 'field list'")
}

func main() {
	flag.Parse()
	if &serviceType == nil {
		fmt.Printf("No service type given. Aborting...\n")
		os.Exit(1)
	}
	if crudFlag && togglerFlag {
		fields = append(fields, "CrudTogglerEndpoint")
	} else {
		if crudFlag {
			fields = append(fields, "CrudEndpoint")
		}
		if togglerFlag {
			fields = append(fields, "TogglerEndpoint")
		}
	}
	if len(fields) == 0 {
		fmt.Printf("No implementing fields given. Aborting...\n")
		os.Exit(1)
	}
	mappedFields := mapFields(fields)
	service := &service{
		Type:     serviceType,
		Param:    strings.ToLower(serviceType),
		Fields:   mappedFields,
		Crud:     crudFlag,
		Toggler:  togglerFlag,
		Scaffold: scaffoldFlag,
	}
	fname := fmt.Sprintf(fileNameTmpl, SnakeCase(serviceType))
	err := writeFile(fname, fileTemplate, service)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fname = fmt.Sprintf(testfileNameTmpl, SnakeCase(serviceType))
	err = writeFile(fname, testfileTemplate, service)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func writeFile(fname string, template *template.Template, service *service) error {
	var buf bytes.Buffer
	err := template.Execute(&buf, service)
	if err != nil {
		fmt.Printf("There was an error parsing the given file template.\n")
		return err
	}
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "", buf.String(), parser.ParseComments)
	if err != nil {
		fmt.Printf("There was an error parsing the given file.\n")
		return err
	}
	file, err := os.Create(fname)
	if err != nil {
		fmt.Printf("There was an error creating the file: %s\n", fname)
		return err
	}
	err = printer.Fprint(file, fileSet, f)
	if err != nil {
		fmt.Printf("There was an error while printing the generated file '%s'.\n", fname)
		return err
	}
	return nil
}

type service struct {
	Type         string
	Param        string
	PayloadType  string
	PayloadParam string
	Fields       map[string]string
	Crud         bool
	Toggler      bool
	Scaffold     bool
}

type field struct {
	Type string
	Name string
}

func SnakeCase(in string) string {
	runes := []rune(in)
	var parts []string
	var lastPos int
	for i := 0; i < len(runes); i++ {
		if i > 0 && unicode.IsUpper(runes[i]) {
			parts = append(parts, strings.ToLower(in[lastPos:i]))
			lastPos = i
		}
	}
	if lastWord := in[lastPos:]; lastWord != "" {
		parts = append(parts, strings.ToLower(lastWord))
	}
	return strings.Join(parts, "_")
}

func typeToName(typ string) string {
	var upcaseIndices []int
	idx := strings.IndexFunc(typ, unicode.IsUpper)
	for idx != -1 {
		upcaseIndices = append(upcaseIndices, idx)
		idx = strings.IndexFunc(typ[idx:], unicode.IsUpper)
	}
	return ""
}

func mapFields(data []string) map[string]string {
	ifaceMap := make(map[string]string)
	for _, iface := range data {
		lastUpcase := strings.LastIndexFunc(iface, unicode.IsUpper)
		r, size := utf8.DecodeRuneInString(iface[lastUpcase:])
		rest := iface[lastUpcase:][size:]
		ifaceMap[iface] = fmt.Sprintf("%c%s", unicode.ToLower(r), rest)
	}
	return ifaceMap
}

type stringsFlag []string

func (v *stringsFlag) Set(s string) error {
	var err error
	*v, err = splitQuotedFields(s)
	if *v == nil {
		*v = []string{}
	}
	return err
}

func splitQuotedFields(s string) ([]string, error) {
	// Split fields allowing '' or "" around elements.
	// Quotes further inside the string do not count.
	var f []string
	for len(s) > 0 {
		for len(s) > 0 && isSpaceByte(s[0]) {
			s = s[1:]
		}
		if len(s) == 0 {
			break
		}
		// Accepted quoted string. No unescaping inside.
		if s[0] == '"' || s[0] == '\'' {
			quote := s[0]
			s = s[1:]
			i := 0
			for i < len(s) && s[i] != quote {
				i++
			}
			if i >= len(s) {
				return nil, fmt.Errorf("unterminated %c string", quote)
			}
			f = append(f, s[:i])
			s = s[i+1:]
			continue
		}
		i := 0
		for i < len(s) && !isSpaceByte(s[i]) {
			i++
		}
		f = append(f, s[:i])
		s = s[i:]
	}
	return f, nil
}

func (v *stringsFlag) String() string {
	return "<stringsFlag>"
}

func isSpaceByte(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

var fileTemplateContent = `{{if not .Scaffold}}// DO NOT EDIT!
// This file is generated by the api generator.{{end}}

// +build !feature

package harvest

import (
	"net/url"
)

type {{.Type}}Service struct {
	{{range $type, $param := .Fields}}
	{{$param}} {{$type}} {{end}}
}

func New{{.Type}}Service({{range $type, $param := .Fields}}{{$param}} {{$type}}, {{end}}) *{{.Type}}Service {
	service := {{.Type}}Service{ {{range $type, $param  := .Fields}}
	{{$param}}: {{$param}},{{end}}
}
return &service
}

{{if .Crud}}{{template "crud" .}}{{end}}
{{if .Toggler}}
func (s *{{.Type}}Service) Toggle({{.Param}} *{{.Type}}) error {
	return s.endpoint.Toggle({{.Param}})
}{{end}}
`

var crudTemplateContent = `
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
`

var testFileContent = `{{if not .Scaffold}}// DO NOT EDIT!
// This file is generated by the api generator.{{end}}

// +build !feature

package harvest

import (
	"net/url"
	"reflect"
	"testing"
)

var (
	expected{{.Type}}ServiceParams = url.Values{"foo": []string{"bar"}}

	tests{{.Type}}Service = map[string]struct { // apiFn to testData
		testData *apiWrapperTestData
		testFn   testFunc
		args     []interface{}
	}{
		{{if .Crud}}"All": {
			&apiWrapperTestData{
				expectedParams:       expected{{.Type}}ServiceParams,
				expectedDataType:     reflect.TypeOf(&[]*{{.Type}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiAllWrapper,
			[]interface{}{&[]*{{.Type}}{}, expected{{.Type}}ServiceParams},
		},
		"Find": {
			&apiWrapperTestData{
				expectedParams:       expected{{.Type}}ServiceParams,
				expectedIdType:       reflect.TypeOf(12),
				expectedDataType:     reflect.TypeOf(&{{.Type}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiFindWrapper,
			[]interface{}{12, &{{.Type}}{}, expected{{.Type}}ServiceParams},
		},
		"Create": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&{{.Type}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiCreateWrapper,
			[]interface{}{&{{.Type}}{}},
		},
		"Update": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&{{.Type}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiUpdateWrapper,
			[]interface{}{&{{.Type}}{}},
		},
		"Delete": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&{{.Type}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiDeleteWrapper,
			[]interface{}{&{{.Type}}{}},
		},{{end}}
		{{if .Toggler}}"Toggle": {
			&apiWrapperTestData{
				expectedDataType:     reflect.TypeOf(&{{.Type}}{}),
				expectedErrorMessage: "ERR",
			},
			testApiToggleWrapper,
			[]interface{}{&{{.Type}}{}},
		},{{end}}
	}
)
{{if .Crud}}
func Test{{.Type}}ServiceAll(t *testing.T) {
	test{{.Type}}ServiceMethod(t, "All")
}

func Test{{.Type}}ServiceFind(t *testing.T) {
	test{{.Type}}ServiceMethod(t, "Find")
}

func Test{{.Type}}ServiceCreate(t *testing.T) {
	test{{.Type}}ServiceMethod(t, "Create")
}

func Test{{.Type}}ServiceUpdate(t *testing.T) {
	test{{.Type}}ServiceMethod(t, "Update")
}

func Test{{.Type}}ServiceDelete(t *testing.T) {
	test{{.Type}}ServiceMethod(t, "Delete")
}{{end}}
{{if .Toggler}}
func Test{{.Type}}ServiceToggle(t *testing.T) {
	test{{.Type}}ServiceMethod(t, "Toggle")
}{{end}}

func test{{.Type}}ServiceMethod(t *testing.T, name string) {
	called := false
	test, ok := tests{{.Type}}Service[name]
	if !ok {
		t.Logf("No test data for method '%s' defined.\n", name)
		t.FailNow()
	}
	api := test.testFn(test.testData, &called)
	service := &{{.Type}}Service{endpoint: api}
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
