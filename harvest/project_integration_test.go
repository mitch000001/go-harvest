// +build integration

package harvest_test

import (
	"testing"
	"time"

	"github.com/mitch000001/go-harvest/harvest"
)

func TestFindAllProjects(t *testing.T) {
	client := createClient(t)
	projects, err := client.Projects.All()
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if projects == nil {
		t.Fatal("Expected projects not to be nil")
	}
	if len(projects) == 0 {
		t.Fatal("Expected projects not to be empty")
	}
	for _, p := range projects {
		t.Logf("Project: %+#v\n", p)
	}
}

func TestFindAllProjectsUpdatedSince(t *testing.T) {
	client := createClient(t)
	projects, err := client.Projects.AllUpdatedSince(time.Now())
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if projects == nil {
		t.Fatal("Expected projects not to be nil")
	}
	if len(projects) != 0 {
		t.Fatalf("Expected projects to have 0 items, got %d", len(projects))
	}
}

func TestFindProject(t *testing.T) {
	client := createClient(t)
	// Find first project
	projects, err := client.Projects.All()
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if projects == nil || len(projects) == 0 {
		t.Fatal("Expected projects not to be nil or empty")
	}
	first := projects[0]

	project, err := client.Projects.Find(first.ID)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if project == nil {
		t.Fatal("Expected project not to be nil")
	}
	if project.ID != first.ID {
		t.Fatalf("Expected to find project with id '%d', got id '%d'\n", first.Id, project.Id)
	}

	// Search unknown id
	project, err = client.Projects.Find(1)
	if err != nil {
		expectedMessage := "No project found for id 1"
		if err.Error() != expectedMessage {
			t.Fatalf("Expected error with message '%s', got error %T with message: %s\n", expectedMessage, err, err.Error())
		}
	}
	if project != nil {
		t.Fatal("Expected project to be nil, got %+#v\n", project)
	}
}

func TestCreateAndDeleteProject(t *testing.T) {
	client := createClient(t)
	project := harvest.Project{
		Name:     "foo",
		ClientId: 2605222,
		BillBy:   "none",
		Budget:   100.00,
		BudgetBy: "none",
	}

	createdProject, err := client.Projects.Create(&project)

	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if createdProject == nil {
		t.Fatal("Expected project not to be nil")
	}
	if createdProject.ID == 0 {
		t.Fatal("Expected Id to be set")
	}
	t.Logf("Got returned project: %+#v\n", createdProject)

	deleted, err := client.Projects.Delete(createdProject)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if !deleted {
		t.Fatalf("Could not delete project created for test")
	}
}

func TestUpdateProject(t *testing.T) {
	client := createClient(t)
	project := harvest.Project{
		Name:     "foo",
		ClientId: 2605222,
		BillBy:   "none",
		Budget:   100.00,
		BudgetBy: "none",
		Billable: true,
	}

	createdProject, err := client.Projects.Create(&project)
	if err != nil {
		panic(err)
	}
	defer client.Projects.Delete(createdProject)

	createdProject.Billable = false

	updatedProject, err := client.Projects.Update(createdProject)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if updatedProject == nil {
		t.Fatal("Expected project not to be nil")
	}
	if updatedProject.Billable != createdProject.Billable {
		t.Fatal("Expected updated project billable to be false, got true")
	}
}

func TestToggleProject(t *testing.T) {
	client := createClient(t)
	project := harvest.Project{
		Name:     "foo",
		ClientId: 2605222,
		BillBy:   "none",
		Budget:   100.00,
		BudgetBy: "none",
		Billable: true,
	}
	createdProject, err := client.Projects.Create(&project)
	if err != nil {
		panic(err)
	}
	defer client.Projects.Delete(createdProject)

	err = client.Projects.Toggle(createdProject)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
}
