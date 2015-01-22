package harvest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

func (c *Client) Type() string {
	return "Client"
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
	return json.Marshal(fmt.Sprintf("%s,%s", tf.StartDate.Format("2006-01-02"), tf.EndDate.Format("2006-01-02")))
}

func (tf *Timeframe) UnmarshalJSON(data []byte) error {
	unquotedData, _ := strconv.Unquote(string(data))
	dates := strings.Split(unquotedData, ",")
	if len(dates) != 2 {
		tf = &Timeframe{}
		return nil
	}
	startTime, _ := time.Parse("2006-01-02", dates[0])
	startDate := ShortDate{startTime}
	endTime, _ := time.Parse("2006-01-02", dates[1])
	endDate := ShortDate{endTime}
	tf = &Timeframe{StartDate: startDate, EndDate: endDate}
	return nil
}

type ClientPayload struct {
	ErrorPayload
	Client *Client `json:"client,omitempty"`
}
