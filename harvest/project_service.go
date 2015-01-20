// +build integration

package harvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

type ProjectService struct {
	api *JsonApi
}

func NewProjectService(api *JsonApi) *ProjectService {
	return &ProjectService{api}
}

func (p *ProjectService) All() ([]*Project, error) {
	response, err := p.api.ProcessRequest("GET", "/projects", nil)
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

func (p *ProjectService) AllUpdatedSince(updatedSince time.Time) ([]*Project, error) {
	params := make(url.Values)
	if !updatedSince.IsZero() {
		params.Add("updated_since", updatedSince.UTC().String())
	}
	query := ""
	if len(params) > 0 {
		query = "?" + params.Encode()
	}
	response, err := p.api.ProcessRequest("GET", fmt.Sprintf("/projects%s", query), nil)
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

func (p *ProjectService) Find(id int) (*Project, error) {
	response, err := p.api.ProcessRequest("GET", fmt.Sprintf("/projects/%d", id), nil)
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

func (p *ProjectService) Create(project *Project) (*Project, error) {
	marshaledProject, err := json.Marshal(ProjectPayload{Project: project})
	if err != nil {
		return nil, err
	}
	response, err := p.api.ProcessRequest("POST", "/projects", bytes.NewReader(marshaledProject))
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

func (p *ProjectService) Delete(project *Project) (bool, error) {
	response, err := p.api.ProcessRequest("DELETE", fmt.Sprintf("/projects/%d", project.Id), nil)
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

func (p *ProjectService) Update(project *Project) (*Project, error) {
	marshaledProject, err := json.Marshal(ProjectPayload{Project: project})
	if err != nil {
		return nil, err
	}
	response, err := p.api.ProcessRequest("PUT", fmt.Sprintf("/projects/%d", project.Id), bytes.NewReader(marshaledProject))
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

func (p *ProjectService) Toggle(project *Project) error {
	response, err := p.api.ProcessRequest("PUT", fmt.Sprintf("/projects/%d/toggle", project.Id), nil)
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
