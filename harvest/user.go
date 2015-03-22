package harvest

import (
	"fmt"
	"time"
)

//go:generate go run ../cmd/api_gen/api_gen.go -type=User -c -t -fields CrudEndpointProvider

type User struct {
	ID                           int       `json:"id,omitempty"`
	Email                        string    `json:"email,omitempty"`
	FirstName                    string    `json:"first_name,omitempty"`
	LastName                     string    `json:"last_name,omitempty"`
	HasAccessToAllFutureProjects bool      `json:"has_access_to_all_future_projects,omitempty"`
	DefaultHourlyRate            float64   `json:"default_hourly_rate,omitempty"`
	IsActive                     bool      `json:"is_active,omitempty"`
	IsAdmin                      bool      `json:"is_admin,omitempty"`
	IsContractor                 bool      `json:"is_contractor,omitempty"`
	Telephone                    string    `json:"telephone,omitempty"`
	Department                   string    `json:"department,omitempty"`
	Timezone                     string    `json:"timezone,omitempty"`
	UpdatedAt                    time.Time `json:"updated_at,omitempty"`
	CreatedAt                    time.Time `json:"created_at,omitempty"`
}

func (u *User) Type() string {
	return "User"
}

func (u *User) Id() int {
	return u.ID
}

func (u *User) SetId(id int) {
	u.ID = id
}

func (u *User) ToggleActive() bool {
	u.IsActive = !u.IsActive
	return u.IsActive
}

type UserPayload struct {
	ErrorPayload
	User *User `json:"user,omitempty"`
}

func (u *UserService) DayEntries(user *User) *DayEntryService {
	id := user.Id()
	userPath := u.endpoint.Path()
	path := fmt.Sprintf("%s/%d/entries", userPath, id)
	endpoint := u.provider.CrudEndpoint(path)
	return NewDayEntryService(endpoint)
}

func (u *UserService) Expenses(user *User) *ExpenseService {
	id := user.Id()
	userPath := u.endpoint.Path()
	path := fmt.Sprintf("%s/%d/expenses", userPath, id)
	endpoint := u.provider.CrudEndpoint(path)
	return NewExpenseService(endpoint)
}
