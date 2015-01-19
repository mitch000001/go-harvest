package harvest

import (
	"encoding/json"
	"fmt"
	"time"
)

//go:generate go run ../cmd/api_gen/api_gen.go -type=Client

type Client struct {
	Name                    string    `json:"name,omitempty"`
	CreatedAt               time.Time `json"created-at,omitempty"`
	UpdatedAt               time.Time `json"updated-at,omitempty"`
	HighriseId              int       `json:"highrise-id,omitempty"`
	ID                      int       `json:"id,omitempty"`
	CacheVersion            int       `json:"cache-version,omitempty"`
	Currency                string    `json:"currency,omitempty"`
	CurrencySymbol          string    `json:"currency-symbol,omitempty"`
	Active                  bool      `json:"active,omitempty"`
	Details                 string    `json:"details,omitempty"`
	DefaultInvoiceTimeframe Timeframe `json:"default-invoice-timeframe,omitempty"`
	LastInvoiceKind         string    `json:"last-invoice-kind,omitempty"`
}

func (c *Client) Id() int {
	return c.ID
}

func (c *Client) SetId(id int) {
	c.ID = id
}

func (c *Client) ToggleActive() bool {
	c.Active = !c.Active
	return c.Active
}

type Timeframe struct {
	StartDate ShortDate
	EndDate   ShortDate
}

func (tf Timeframe) MarshalJSON() ([]byte, error) {
	if tf.StartDate.IsZero() || tf.EndDate.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(fmt.Sprintf("%s,%s", time.Time(tf.StartDate).Format("2006-01-02"), time.Time(tf.EndDate).Format("2006-01-02")))
}

func (tf *Timeframe) UnmarshalJSON(data []byte) error {
	strDate := string(data)
	var startDateString string
	var endDateString string
	_, err := fmt.Sscanf(strDate, "\"%s,%s\"", &startDateString, &endDateString)
	if err != nil {
		tf = &Timeframe{}
		return nil
	}
	var startDate ShortDate
	var endDate ShortDate
	startTime, err := time.Parse("2006-01-02", startDateString)
	if err != nil {
		startDate = ShortDate{}
		err = nil
	} else {
		startDate = ShortDate(startTime)
	}
	endTime, err := time.Parse("2006-01-02", startDateString)
	if err != nil {
		endDate = ShortDate{}
		err = nil
	} else {
		endDate = ShortDate(endTime)
	}
	if startDate.IsZero() || endDate.IsZero() {
		tf = &Timeframe{}
	} else {
		tf = &Timeframe{StartDate: startDate, EndDate: endDate}
	}
	return nil
}

type ClientPayload struct {
	ErrorPayload
	Client *Client `json:"client,omitempty"`
}
