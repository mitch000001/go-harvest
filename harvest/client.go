package harvest

import "time"

//go:generate go run ../cmd/api_gen/api_gen.go -type=Client -c -t

type Client struct {
	Name                    string    `json:"name,omitempty"`
	CreatedAt               time.Time `json:"created-at,omitempty"`
	UpdatedAt               time.Time `json:"updated-at,omitempty"`
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

type ClientPayload struct {
	ErrorPayload
	Client *Client `json:"client,omitempty"`
}
