package harvest

import "testing"

func TestUserToggleActive(t *testing.T) {
	user := &User{
		IsActive: true,
	}
	status := user.ToggleActive()

	if status {
		t.Logf("Expected status to be false, got true\n")
		t.Fail()
	}

	if user.IsActive {
		t.Logf("Expected IsActive to be false, got true\n")
		t.Fail()
	}
}
