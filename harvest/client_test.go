package harvest

import "testing"

func TestClientToggleActive(t *testing.T) {
	client := &Client{
		Active: true,
	}
	status := client.ToggleActive()

	if status {
		t.Logf("Expected status to be false, got true\n")
		t.Fail()
	}

	if client.Active {
		t.Logf("Expected IsActive to be false, got true\n")
		t.Fail()
	}
}
