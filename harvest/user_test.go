package harvest

import "testing"

func TestUserSetId(t *testing.T) {
	user := &User{}

	if user.ID != 0 {
		t.Logf("Expected id to be 0, got %d\n", user.ID)
		t.Fail()
	}

	user.SetId(12)

	if user.ID != 12 {
		t.Logf("Expected id to be 12, got %d\n", user.ID)
		t.Fail()
	}
}

func TestUserId(t *testing.T) {
	user := &User{}

	if user.Id() != 0 {
		t.Logf("Expected id to be 0, got %d\n", user.ID)
		t.Fail()
	}

	user.ID = 12

	if user.Id() != 12 {
		t.Logf("Expected id to be 12, got %d\n", user.ID)
		t.Fail()
	}
}

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

func TestUserType(t *testing.T) {
	typ := (&User{}).Type()

	if typ != "User" {
		t.Logf("Expected Type to equal 'User', got '%s'\n", typ)
		t.Fail()
	}
}
