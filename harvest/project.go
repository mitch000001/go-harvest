package harvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

type ShortDate time.Time

func (date ShortDate) IsZero() bool {
	return time.Time(date).IsZero()
}

func (date ShortDate) MarshalJSON() ([]byte, error) {
	if date.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(time.Time(date).Format("2006-01-02"))
}

func (date *ShortDate) UnmarshalJSON(data []byte) error {
	strDate := string(data)
	time, err := time.Parse("2006-01-02", strDate[1:len(strDate)-1])
	if err != nil {
		date = &ShortDate{}
		err = nil
	} else {
		*date = ShortDate(time)
	}
	return err
}

//go:generate go run ../cmd/api_gen/api_gen.go -type=Project

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

type ProjectsService struct {
	h *Harvest
}

func NewProjectsService(client *Harvest) *ProjectsService {
	return &ProjectsService{client}
}

func (p *ProjectsService) All() ([]*Project, error) {
	response, err := p.h.ProcessRequest("GET", "/projects", nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	projectResponses := make([]*ProjectPayload, 0)
	err = json.Unmarshal(responseBytes, &projectResponses)
	if err != nil {
		return nil, err
	}
	projects := make([]*Project, len(projectResponses))
	for i, p := range projectResponses {
		projects[i] = p.Project
	}
	return projects, nil
}

func (p *ProjectsService) AllUpdatedSince(updatedSince time.Time) ([]*Project, error) {
	params := make(url.Values)
	if !updatedSince.IsZero() {
		params.Add("updated_since", updatedSince.UTC().String())
	}
	query := ""
	if len(params) > 0 {
		query = "?" + params.Encode()
	}
	response, err := p.h.ProcessRequest("GET", fmt.Sprintf("/projects%s", query), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	projectResponses := make([]*ProjectPayload, 0)
	err = json.Unmarshal(responseBytes, &projectResponses)
	if err != nil {
		return nil, err
	}
	projects := make([]*Project, len(projectResponses))
	for i, p := range projectResponses {
		projects[i] = p.Project
	}
	return projects, nil
}

func (p *ProjectsService) Find(id int) (*Project, error) {
	response, err := p.h.ProcessRequest("GET", fmt.Sprintf("/projects/%d", id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == 404 {
		return nil, &ResponseError{&ErrorPayload{fmt.Sprintf("No project found for id %d", id)}}
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	projectResponse := ProjectPayload{}
	err = json.Unmarshal(responseBytes, &projectResponse)
	if err != nil {
		return nil, err
	}
	return projectResponse.Project, nil
}

func (p *ProjectsService) Create(project *Project) (*Project, error) {
	marshaledProject, err := json.Marshal(ProjectPayload{Project: project})
	if err != nil {
		return nil, err
	}
	response, err := p.h.ProcessRequest("POST", "/projects", bytes.NewReader(marshaledProject))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == 201 {
		location := response.Header.Get("Location")
		projectId := -1
		fmt.Sscanf(location, "/projects/%d", &projectId)
		if projectId == -1 {
			return nil, &ResponseError{&ErrorPayload{"No id for project received"}}
		}
		project.SetId(projectId)
		return project, nil
	}
	if response.StatusCode == 200 {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			return nil, err
		}
		return nil, &ResponseError{&apiResponse}
	}
	return nil, &ResponseError{&ErrorPayload{response.Status}}
}

func (p *ProjectsService) Delete(project *Project) (bool, error) {
	response, err := p.h.ProcessRequest("DELETE", fmt.Sprintf("/projects/%d", project.Id), nil)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return true, nil
	} else if response.StatusCode == 400 {
		return false, nil
	} else {
		panic(response.Status)
	}
}

func (p *ProjectsService) Update(project *Project) (*Project, error) {
	marshaledProject, err := json.Marshal(ProjectPayload{Project: project})
	if err != nil {
		return nil, err
	}
	response, err := p.h.ProcessRequest("PUT", fmt.Sprintf("/projects/%d", project.Id), bytes.NewReader(marshaledProject))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		apiResponse := ErrorPayload{}
		err = json.Unmarshal(responseBytes, &apiResponse)
		if err != nil {
			return nil, err
		}
		return nil, &ResponseError{&apiResponse}
	}
	return project, nil
}

func (p *ProjectsService) Toggle(project *Project) error {
	response, err := p.h.ProcessRequest("PUT", fmt.Sprintf("/projects/%d/toggle", project.Id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return nil
	} else if response.StatusCode == 400 {
		hint := response.Header.Get("Hint")
		return &ResponseError{&ErrorPayload{hint}}
	} else {
		panic(response.Status)
	}
}
