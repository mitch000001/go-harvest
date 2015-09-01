// +build !feature

package harvest

import (
	"net/url"
)

type TaskAssignmentService struct {
	endpoint CrudEndpoint
}

func NewTaskAssignmentService(endpoint CrudEndpoint) *TaskAssignmentService {
	service := TaskAssignmentService{
		endpoint: endpoint,
	}
	return &service
}

func (s *TaskAssignmentService) All(taskassignments *[]*TaskAssignment, params url.Values) error {
	return s.endpoint.All(taskassignments, params)
}

func (s *TaskAssignmentService) Find(id int, taskassignment *TaskAssignment, params url.Values) error {
	return s.endpoint.Find(id, taskassignment, params)
}

func (s *TaskAssignmentService) Create(taskassignment *TaskAssignment) error {
	return s.endpoint.Create(taskassignment)
}

func (s *TaskAssignmentService) Update(taskassignment *TaskAssignment) error {
	return s.endpoint.Update(taskassignment)
}

func (s *TaskAssignmentService) Delete(taskassignment *TaskAssignment) error {
	return s.endpoint.Delete(taskassignment)
}
