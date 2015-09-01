package harvest

import "time"

//go:generate go run ../cmd/api_gen/api_gen.go -type=TaskAssignment -c -s

type TaskAssignment struct {
	ID        int `json:"id"`
	TaskId    int `json:"task-id"`
	ProjectId int `json:"project-id"`
	// True if task is billable for this project
	Billable bool `json:"billable"`
	// True if task was deactivated for project, preventing further hours to be logged against it
	Deactivated bool `json:"deactivated"`
	// The budget (if present) for the task in project
	Budget float64 `json:"budget"`
	// The hourly rate (if present) for the task in project
	HourlyRate float64   `json:"hourly-rate"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

func (t *TaskAssignment) Id() int {
	return t.ID
}

func (t *TaskAssignment) SetId(id int) {
	t.ID = id
}

func (t *TaskAssignment) Type() string {
	return "task-assignment"
}
