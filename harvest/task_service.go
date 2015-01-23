package harvest

import (
	"fmt"
	"net/http"
)

func (t *TaskService) Activate(task *Task) error {
	response, err := t.processor.Process("POST", fmt.Sprintf("/tasks/%d", task.Id()), nil)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusOK {
		task.Deactivated = false
	}
	return nil
}
