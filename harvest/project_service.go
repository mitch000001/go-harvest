package harvest

import "fmt"

func (p *ProjectService) UserAssignments(project *Project) *UserAssignmentService {
	id := project.Id()
	projectPath := p.endpoint.Path()
	path := fmt.Sprintf("%s/%d/user_assignments", projectPath, id)
	endpoint := p.provider.CrudEndpoint(path)
	return NewUserAssignmentService(endpoint)
}
