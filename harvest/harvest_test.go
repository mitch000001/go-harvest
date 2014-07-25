package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseSubdomain(t *testing.T) {
	// only the subdomain name given
	subdomain := "foo"

	testSubdomain(subdomain, t)

	subdomain = "https://foo.harvestapp.com/"

	testSubdomain(subdomain, t)
}

func testSubdomain(subdomain string, t *testing.T) {
	testUrl, err := parseSubdomain(subdomain)
	if err != nil {
		t.Fatal(err)
	}
	if testUrl == nil {
		t.Fatal("Expected url not to be nil")
	}
	if testUrl.String() != "https://foo.harvestapp.com/" {
		t.Fatalf("Expected schema to equal 'https://foo.harvestapp.com/', got '%s'", testUrl)
	}
}

func createClient(t *testing.T) *Harvest {
	subdomain := os.Getenv("HARVEST_SUBDOMAIN")
	username := os.Getenv("HARVEST_USERNAME")
	password := os.Getenv("HARVEST_PASSWORD")

	client, err := NewBasicAuthClient(subdomain, &BasicAuthConfig{username, password})
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("Expected client not to be nil")
	}
	return client
}

func testAllFunc(allFunc interface{}, options interface{}, t *testing.T) {
	allFuncVal := reflect.ValueOf(allFunc)
	optionsValue := reflect.ValueOf(options)
	args := []reflect.Value{}
	if optionsValue.IsValid() {
		args = append(args, optionsValue.Elem())
	} else {
		allFuncType := allFuncVal.Type()
		argsType := allFuncType.In(0)
		args = append(args, reflect.Zero(argsType))
	}
	result := allFuncVal.Call(args)
	payload := result[0]
	err := result[1]
	if !err.IsNil() {
		errorMessage := err.MethodByName("Error").Call([]reflect.Value{})[0]
		t.Fatalf("Got error %T with message: %s\n", err.Type, errorMessage)
	}
	if payload.Len() != 1 {
		t.Fatalf("Expected 1 user, got %d", payload.Len())
	}
	if payload.Index(0).IsNil() {
		t.Fatal("Expected payload not to be nil")
	}
	for i := 0; i < payload.Len(); i++ {
		p := payload.Index(i)
		t.Logf("%s: %+#v\n", p.Type(), p.Interface())
	}
}

func testFindFunc(findFunc interface{}, id int, t *testing.T) {
	findFuncVal := reflect.ValueOf(findFunc)
	// Find existing entity
	result := findFuncVal.Call([]reflect.Value{reflect.ValueOf(id)})
	payload := result[0]
	err := result[1]
	if !err.IsNil() {
		errorMessage := err.MethodByName("Error").Call([]reflect.Value{})[0]
		t.Fatalf("Got error %T with message: %s\n", err.Type, errorMessage)
	}
	if payload.IsNil() {
		t.Fatalf("Expect to find entity with id %d, got nil\n", id)
	}
	payloadId := payload.Elem().FieldByName("Id").Int()
	if payloadId != int64(id) {
		t.Fatalf("Expect to find entity with id %d, got %#v\n", id, payload.Interface())
	}

	// No entity with that id
	result = findFuncVal.Call([]reflect.Value{reflect.ValueOf(1)})
	payload = result[0]
	err = result[1]
	if !err.IsNil() {
		expectedErrorMessage := "Id not found: 1"
		errorMessage := err.MethodByName("Error").Call([]reflect.Value{})[0].String()
		if errorMessage != expectedErrorMessage {
			t.Fatalf("Expected ResponseError with message '%s', got error %s with message: %s\n", expectedErrorMessage, err.Type(), errorMessage)
		}
	}
	if !payload.IsNil() {
		t.Fatalf("Expected entity to be nil, got '%+#v'", payload.Interface())
	}
}

func testCreateAndDeleteFunc(createFunc interface{}, deleteFunc interface{}, entityToCreate interface{}, t *testing.T) {
	createFuncVal := reflect.ValueOf(createFunc)
	result := createFuncVal.Call([]reflect.Value{reflect.ValueOf(entityToCreate)})
	payload := result[0]
	err := result[1]
	if !err.IsNil() {
		errorMessage := err.MethodByName("Error").Call([]reflect.Value{})[0]
		t.Fatalf("Got error %T with message: %s\n", err.Type, errorMessage)
	}
	if payload.IsNil() {
		t.Fatal("Expect entity not to be nil\n")
	}
	t.Logf("Got returned %s: %+#v\n", payload.Type(), payload.Interface())
	deleteFuncVal := reflect.ValueOf(deleteFunc)
	result = deleteFuncVal.Call([]reflect.Value{payload})
	err = result[0]
	if !err.IsNil() {
		errorMessage := err.MethodByName("Error").Call([]reflect.Value{})[0]
		t.Fatalf("Got error %T with message: %s\n", err.Type, errorMessage)
	}
}

func testUpdateFunc(updateFunc interface{}, entityToUpdate interface{}, fieldToUpdate string, t *testing.T) {
	entityVal := reflect.ValueOf(entityToUpdate)

	updateFuncVal := reflect.ValueOf(updateFunc)
	result := updateFuncVal.Call([]reflect.Value{entityVal})
	payload := result[0]
	err := result[1]
	if !err.IsNil() {
		errorMessage := err.MethodByName("Error").Call([]reflect.Value{})[0]
		t.Fatalf("Got error %T with message: %s\n", err.Type, errorMessage)
	}
	if payload.IsNil() {
		t.Fatal("Expect entity not to be nil\n")
	}
	t.Logf("Got returned %s: %+#v\n", payload.Type(), payload.Interface())
	oldVal := entityVal.Elem().FieldByName(fieldToUpdate).Interface()
	newVal := payload.Elem().FieldByName(fieldToUpdate).Interface()
	if oldVal != newVal {
		t.Fatalf("Expected updated field %s to equal '%s', got '%s'", fieldToUpdate, oldVal, newVal)
	}
}

func testUpdateFuncInvalidUpdate(updateFunc interface{}, entityToUpdate interface{}, t *testing.T) {
	entityVal := reflect.ValueOf(entityToUpdate)

	updateFuncVal := reflect.ValueOf(updateFunc)
	result := updateFuncVal.Call([]reflect.Value{entityVal})
	payload := result[0]
	err := result[1]
	if err.IsNil() {
		t.Fatal("Expected ResponseError, got nil")
	}
	if !payload.IsNil() {
		t.Fatalf("Expect entity to be nil, got '%+#v'\n", payload.Interface())
	}
}
