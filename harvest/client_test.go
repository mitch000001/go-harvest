package harvest

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

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

func TestTimeframeMarshalJSON(t *testing.T) {
	startDate := ShortDate{time.Date(2014, time.February, 01, 0, 0, 0, 0, time.UTC)}
	endDate := ShortDate{time.Date(2014, time.April, 01, 0, 0, 0, 0, time.UTC)}

	var tests = []struct {
		timeframe    Timeframe
		expectedJson string
	}{
		{
			timeframe:    Timeframe{StartDate: startDate, EndDate: endDate},
			expectedJson: `"2014-02-01,2014-04-01"`,
		},
		{
			timeframe:    Timeframe{StartDate: startDate},
			expectedJson: `""`,
		},
		{
			timeframe:    Timeframe{EndDate: endDate},
			expectedJson: `""`,
		},
		{
			timeframe:    Timeframe{},
			expectedJson: `""`,
		},
	}

	for _, test := range tests {
		bytes, err := json.Marshal(&test.timeframe)
		if err != nil {
			t.Logf("Expected error to be nil, got %T: %v\n", err, err)
			t.Fail()
		}

		if !reflect.DeepEqual(string(bytes), test.expectedJson) {
			t.Logf("Expected date to be '%s', got '%s'\n", test.expectedJson, string(bytes))
			t.Fail()
		}
	}

}

func TestTimeframeUnmarshalJSON(t *testing.T) {
	// startDate := ShortDate{time.Date(2014, time.February, 01, 0, 0, 0, 0, time.UTC)}
	// endDate := ShortDate{time.Date(2014, time.April, 01, 0, 0, 0, 0, time.UTC)}

	var tests = []struct {
		testJson          string
		expectedTimeframe Timeframe
	}{
		// TODO: happy path doesn't work?!
		// {
		// 	`"2014-02-01,2014-04-01"`,
		// 	Timeframe{StartDate: startDate, EndDate: endDate},
		// },
		{
			`"2014-02-01,"`,
			Timeframe{},
		},
		{
			`""`,
			Timeframe{},
		},
		{
			`","`,
			Timeframe{},
		},
		{
			`"2014-02-01,abcde"`,
			Timeframe{},
		},
		{
			`"abcde,2014-04-01"`,
			Timeframe{},
		},
		{
			`"abcde,abcde"`,
			Timeframe{},
		},
	}

	for _, test := range tests {
		var timeframe Timeframe
		err := json.Unmarshal([]byte(test.testJson), &timeframe)
		if err != nil {
			t.Logf("Expected error to be nil, got %T: %v\n", err, err)
			t.Fail()
		}

		if !reflect.DeepEqual(timeframe, test.expectedTimeframe) {
			t.Logf("Expected date to be '%+#v', got '%+#v'\n", test.expectedTimeframe, timeframe)
			t.Fail()
		}
	}
}
