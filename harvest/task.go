package harvest

import "time"

//go:generate go run ../cmd/api_gen/api_gen.go -type=Task -fields RequestProcessor

type Task struct {
	// If true task will be added as billable upon assigning it to a project
	BillableByDefault bool `json: "billable-by-default"`
	// False if hours can be recorded against this task.  True if task is archived -->
	Deactivated       bool `json:"deactivated"`
	DefaultHourlyRate int  `json: "default-hourly-rate"`
	ID                int  `json:"id"`
	// If true task is added to new projects by default -->
	IsDefault bool      `json: "is-default"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json"created-at,omitempty"`
	UpdatedAt time.Time `json"updated-at,omitempty"`
}

func (t *Task) Type() string {
	return "Task"
}

func (t *Task) Id() int {
	return t.ID
}

func (t *Task) SetId(id int) {
	t.ID = id
}
