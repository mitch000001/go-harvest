package harvest

import "time"

//go:generate go run ../cmd/api_gen/api_gen.go -type=UserAssignment -c -s

type UserAssignment struct {
	ID        int `json:"id"`
	UserId    int `json:"user-id"`
	ProjectId int `json:"project-id"`
	// If true, user cannot log more hours toward the project -->
	Deactivated bool `json:"deactivated"`
	// Hourly rate of user on current project -->
	HourlyRate float64   `json:"hourly-rate"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

func (u *UserAssignment) Id() int {
	return u.ID
}

func (u *UserAssignment) SetId(id int) {
	u.ID = id
}

func (u *UserAssignment) Type() string {
	return "user-assignment"
}
