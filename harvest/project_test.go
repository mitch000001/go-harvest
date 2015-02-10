package harvest

import "testing"

func TestProjectSetId(t *testing.T) {
	project := &Project{}

	if project.ID != 0 {
		t.Logf("Expected id to be 0, got %d\n", project.ID)
		t.Fail()
	}

	project.SetId(12)

	if project.ID != 12 {
		t.Logf("Expected id to be 12, got %d\n", project.ID)
		t.Fail()
	}
}

func TestProjectId(t *testing.T) {
	project := &Project{}

	if project.Id() != 0 {
		t.Logf("Expected id to be 0, got %d\n", project.ID)
		t.Fail()
	}

	project.ID = 12

	if project.Id() != 12 {
		t.Logf("Expected id to be 12, got %d\n", project.ID)
		t.Fail()
	}
}

func TestProjectToggleActive(t *testing.T) {
	project := &Project{
		Active: true,
	}
	status := project.ToggleActive()

	if status {
		t.Logf("Expected status to be false, got true\n")
		t.Fail()
	}

	if project.Active {
		t.Logf("Expected IsActive to be false, got true\n")
		t.Fail()
	}
}

func TestProjectType(t *testing.T) {
	typ := (&Project{}).Type()

	if typ != "Project" {
		t.Logf("Expected Type to equal 'Project', got '%s'\n", typ)
		t.Fail()
	}
}
