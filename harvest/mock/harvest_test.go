package mock

import (
	"reflect"
	"testing"

	"github.com/mitch000001/go-harvest/harvest"
)

func TestNew(t *testing.T) {
	mock := Mock{}

	client, err := New(mock)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if client == nil {
		t.Logf("Expected client not to be nil\n")
		t.FailNow()
	}
}

func TestMockProjects(t *testing.T) {
	t.Skip()
	mock := Mock{}

	client, err := New(mock)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expectedProjects := []*harvest.Project{
		&harvest.Project{},
	}

	var actualProjects []*harvest.Project

	err = client.Projects.All(&actualProjects, nil)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedProjects, actualProjects) {
		t.Logf("Expected projects to equal\n%#v\n\tgot\n%#v\n", expectedProjects, actualProjects)
		t.Fail()
	}
}
