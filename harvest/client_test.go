package main

import (
	"testing"
)

func TestFindAllClients(t *testing.T) {
	client := createClient(t)
	testAllFunc(client.Clients.All, nil, t)
}

func TestFindClient(t *testing.T) {
	t.Skip()
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

	project, err := client.Projects.Find(first.Id)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if project == nil {
		t.Fatal("Expected project not to be nil")
	}
	if project.Id != first.Id {
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
