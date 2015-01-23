package harvest

import "testing"

func TestTaskType(t *testing.T) {
	typ := (&Task{}).Type()

	if typ != "Task" {
		t.Logf("Expected Type to equal 'Task', got '%s'\n", typ)
		t.Fail()
	}
}

func TestTaskSetId(t *testing.T) {
	task := &Task{}

	if task.ID != 0 {
		t.Logf("Expected id to be 0, got %d\n", task.ID)
		t.Fail()
	}

	task.SetId(12)

	if task.ID != 12 {
		t.Logf("Expected id to be 12, got %d\n", task.ID)
		t.Fail()
	}
}

func TestTaskId(t *testing.T) {
	task := &Task{}

	if task.Id() != 0 {
		t.Logf("Expected id to be 0, got %d\n", task.ID)
		t.Fail()
	}

	task.ID = 12

	if task.Id() != 12 {
		t.Logf("Expected id to be 12, got %d\n", task.ID)
		t.Fail()
	}
}
