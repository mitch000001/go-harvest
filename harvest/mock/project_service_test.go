package mock

import (
	"reflect"
	"testing"
	"time"

	"github.com/mitch000001/go-harvest/harvest"
)

func TestProjectServiceAll(t *testing.T) {
	service := ProjectService{
		projects: []*harvest.Project{
			&harvest.Project{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
			&harvest.Project{ID: 2, UpdatedAt: time.Date(2015, 8, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	var actualProjects []*harvest.Project

	err := service.All(&actualProjects, nil)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedProjects := []*harvest.Project{
		&harvest.Project{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		&harvest.Project{ID: 2, UpdatedAt: time.Date(2015, 8, 1, 0, 0, 0, 0, time.UTC)},
	}

	if !reflect.DeepEqual(expectedProjects, actualProjects) {
		t.Logf("Expected projects to equal\n%q\n\tgot\n%q\n", expectedProjects, actualProjects)
		t.Fail()
	}

	// With updated since
	var params harvest.Params
	err = service.All(&actualProjects, *params.UpdatedSince(time.Date(2015, 7, 1, 0, 0, 0, 0, time.UTC)))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedProjects = []*harvest.Project{
		&harvest.Project{ID: 2, UpdatedAt: time.Date(2015, 8, 1, 0, 0, 0, 0, time.UTC)},
	}

	if !reflect.DeepEqual(expectedProjects, actualProjects) {
		t.Logf("Expected projects to equal\n%q\n\tgot\n%q\n", expectedProjects, actualProjects)
		t.Fail()
	}
}
