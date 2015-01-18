package harvest

import "testing"

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
