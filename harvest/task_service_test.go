package harvest

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestTaskServiceActivate(t *testing.T) {
	processor := &testProcessor{}
	taskService := &TaskService{processor: processor}

	task := Task{
		ID:          12,
		Deactivated: true}

	processor.response = &http.Response{
		StatusCode: http.StatusOK,
	}

	err := taskService.Activate(&task)

	if err != nil {
		t.Logf("Expected no error, got %T: %s\n", err, err.Error())
		t.Fail()
	}

	if processor.method != "POST" {
		t.Logf("Expected request method to equal 'POST', got '%q'\n", processor.method)
		t.Fail()
	}

	if processor.path != "/tasks/12" {
		t.Logf("Expected request method to equal '/tasks/12', got '%q'\n", processor.path)
		t.Fail()
	}

	if processor.body != nil {
		t.Logf("Expected request body to be nil, got '%+#v'\n", processor.body)
		t.Fail()
	}

	if task.Deactivated {
		t.Logf("Expected task to be active, was not\n")
		t.Fail()
	}

	// error case
	task.Deactivated = true
	processor.err = fmt.Errorf("ERROR")

	err = taskService.Activate(&task)

	if err == nil {
		t.Logf("Expected error, got nil\n")
		t.Fail()
	}

	if !task.Deactivated {
		t.Logf("Expected task to be deactived, was not\n")
		t.Fail()
	}
}

type testProcessor struct {
	method   string
	path     string
	body     io.Reader
	response *http.Response
	err      error
}

func (t *testProcessor) Process(method string, path string, body io.Reader) (*http.Response, error) {
	t.method = method
	t.path = path
	t.body = body
	return t.response, t.err
}
