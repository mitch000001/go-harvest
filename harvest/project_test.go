package harvest

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

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

func TestShortDateUnmarshalJSON(t *testing.T) {
	testJson := `"2014-02-01"`

	var date ShortDate

	err := json.Unmarshal([]byte(testJson), &date)

	if err != nil {
		t.Logf("Expected error to be nil, got %T: %v\n", err, err)
		t.Fail()
	}

	if &date == nil {
		t.Logf("Expected date not to be nil\n")
		t.Fail()
	}

	expectedDate, err := time.Parse("2006-01-02", "2014-02-01")
	expectedShortDate := ShortDate{expectedDate}

	if err != nil {
		t.Logf("Expected error to be nil, got %T: %v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedShortDate, date) {
		t.Logf("Expected date to be '%+#v', got '%+#v'\n", expectedShortDate, date)
		t.Fail()
	}
}

func TestShortDateMarshalJSON(t *testing.T) {
	date := ShortDate{time.Date(2014, time.February, 01, 0, 0, 0, 0, time.UTC)}

	bytes, err := json.Marshal(&date)

	if err != nil {
		t.Logf("Expected error to be nil, got %T: %v\n", err, err)
		t.Fail()
	}

	expectedJson := `"2014-02-01"`

	if !reflect.DeepEqual(string(bytes), expectedJson) {
		t.Logf("Expected date to be '%s', got '%s'\n", expectedJson, string(bytes))
		t.Fail()
	}

	// Date is Zero
	date = ShortDate{}

	bytes, err = json.Marshal(&date)

	if err != nil {
		t.Logf("Expected error to be nil, got %T: %v\n", err, err)
		t.Fail()
	}

	expectedJson = `""`

	if !reflect.DeepEqual(string(bytes), expectedJson) {
		t.Logf("Expected date to be '%s', got '%s'\n", expectedJson, string(bytes))
		t.Fail()
	}
}
