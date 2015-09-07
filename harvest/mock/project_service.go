package mock

import (
	"fmt"
	"time"

	"github.com/mitch000001/go-harvest/harvest"
)

type ProjectService struct {
	projects []*harvest.Project
}

func (p ProjectService) All(projects *[]*harvest.Project, params harvest.Params) error {
	if params != nil {
		if updatedSince := params.Get("updated_since"); updatedSince != "" {
			t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", updatedSince)
			if err != nil {
				return fmt.Errorf("Error while parsing updated since: %v", err)
			}
			*projects = make([]*harvest.Project, 0)
			for _, p := range p.projects {
				if p.UpdatedAt.After(t) {
					*projects = append(*projects, p)
				}
			}

		}
	} else {
		*projects = p.projects
	}
	return nil
}
