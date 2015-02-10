package harvest

import "testing"

func TestClientSetId(t *testing.T) {
	client := &Client{}

	if client.ID != 0 {
		t.Logf("Expected id to be 0, got %d\n", client.ID)
		t.Fail()
	}

	client.SetId(12)

	if client.ID != 12 {
		t.Logf("Expected id to be 12, got %d\n", client.ID)
		t.Fail()
	}
}

func TestClientId(t *testing.T) {
	client := &Client{}

	if client.Id() != 0 {
		t.Logf("Expected id to be 0, got %d\n", client.ID)
		t.Fail()
	}

	client.ID = 12

	if client.Id() != 12 {
		t.Logf("Expected id to be 12, got %d\n", client.ID)
		t.Fail()
	}
}

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

func TestClientType(t *testing.T) {
	typ := (&Client{}).Type()

	if typ != "Client" {
		t.Logf("Expected Type to equal 'Client', got '%s'\n", typ)
		t.Fail()
	}
}
