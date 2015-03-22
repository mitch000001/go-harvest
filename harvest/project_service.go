package harvest

import "fmt"

func (p *ProjectService) UserAssignments(project *Project) *UserAssignmentService {
	id := project.Id()
	projectPath := p.endpoint.Path()
	path := fmt.Sprintf("%s/%d/user_assignments", projectPath, id)
	endpoint := p.provider.CrudEndpoint(path)
	return NewUserAssignmentService(endpoint)
}

func (p *ProjectService) DayEntries(project *Project) *DayEntryService {
	id := project.Id()
	projectPath := p.endpoint.Path()
	path := fmt.Sprintf("%s/%d/entries", projectPath, id)
	endpoint := p.provider.CrudEndpoint(path)
	return NewDayEntryService(endpoint)
}

func (p *ProjectService) Expenses(project *Project) *ExpenseService {
	id := project.Id()
	projectPath := p.endpoint.Path()
	path := fmt.Sprintf("%s/%d/expenses", projectPath, id)
	endpoint := p.provider.CrudEndpoint(path)
	return NewExpenseService(endpoint)
}
