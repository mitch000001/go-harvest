package harvest

import "time"

//go:generate go run ../cmd/api_gen/api_gen.go -type=Project -c -t -fields CrudEndpointProvider

type Project struct {
	Name     string `json:"name,omitempty"`
	ID       int    `json:"id,omitempty"`
	ClientId int    `json:"client_id,omitempty"`
	Code     string `json:"code,omitempty"`
	Active   bool   `json:"active,omitempty"`
	Notes    string `json:"notes,omitempty"`
	Billable bool   `json:"billable,omitempty"`
	/* Shows if the project is billed by task hourly rate or
	person hourly rate. Options: Tasks, People, none */
	BillBy                    string  `json:"bill_by,omitempty"`
	CostBudget                float64 `json:"cost_budget,omitempty"`
	CostBudgetIncludeExpenses bool    `json:"cost_budget_include_expenses,omitempty"`
	HourlyRate                string  `json:"hourly_rate,omitempty"`
	/* Shows if the budget provided by total project hours,
	total project cost, by tasks, by people or none provided.
	Options: project, project_cost, task, person, none */
	BudgetBy                         string    `json:"budget_by,omitempty"`
	Budget                           float64   `json:"budget,omitempty"`
	NotifyWhenOverBudget             bool      `json:"notify_when_over_budget,omitempty"`
	OverBudgetNotificationPercentage float32   `json:"over_budget_notification_percentage,omitempty"`
	OverBudgetNotifiedAt             string    `json:"over_budget_notified_at,omitempty"`
	ShowBudgetToAll                  bool      `json:"show_budget_to_all,omitempty"`
	CreatedAt                        time.Time `json:"created_at,omitempty"`
	UpdatedAt                        time.Time `json:"updated_at,omitempty"`
	/* These are hints to when the earliest and latest date when a
	timesheet record or an expense was created for a project. Note
	that these fields are only updated once every 24 hours, they
	are useful to constructing a full project timeline. */
	HintEarliestRecordAt ShortDate `json:"hint_earliest_record_at,omitempty"`
	HintLatestRecordAt   ShortDate `json:"hint_latest_record_at,omitempty"`
}

func (p *Project) Type() string {
	return "Project"
}

func (p *Project) Id() int {
	return p.ID
}

func (p *Project) SetId(id int) {
	p.ID = id
}

func (p *Project) ToggleActive() bool {
	p.Active = !p.Active
	return p.Active
}

type ProjectPayload struct {
	ErrorPayload
	Project *Project `json:"project,omitempty"`
}
