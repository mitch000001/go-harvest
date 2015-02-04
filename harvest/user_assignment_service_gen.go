// +build !feature

package harvest

import (
	"net/url"
)

type UserAssignmentService struct {
	endpoint CrudEndpoint
}

func NewUserAssignmentService(endpoint CrudEndpoint) *UserAssignmentService {
	service := UserAssignmentService{
		endpoint: endpoint,
	}
	return &service
}

func (s *UserAssignmentService) All(userassignments *[]*UserAssignment, params url.Values) error {
	return s.endpoint.All(userassignments, params)
}

func (s *UserAssignmentService) Find(id int, userassignment *UserAssignment, params url.Values) error {
	return s.endpoint.Find(id, userassignment, params)
}

func (s *UserAssignmentService) Create(userassignment *UserAssignment) error {
	return s.endpoint.Create(userassignment)
}

func (s *UserAssignmentService) Update(userassignment *UserAssignment) error {
	return s.endpoint.Update(userassignment)
}

func (s *UserAssignmentService) Delete(userassignment *UserAssignment) error {
	return s.endpoint.Delete(userassignment)
}
