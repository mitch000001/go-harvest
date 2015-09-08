package harvest

import (
	"net/url"
	"time"
)

// Params is an enhanced url.Values which can also be used without prior
// initialization. It also defines helper methods to use it as a query params
// builder
//
// The method set to url.Values is identical, but adds lazy initialization.
type Params url.Values

// init initializes the Params type if it's nil
//
// Note that the Params type is a url.Values which itself is a
// map[string][]string, so an uninitialized Params object is in fact an
// uninitialized map, so we call make(Params) to initialize it
func (p *Params) init() {
	if *p == nil {
		*p = make(Params)
	}
}

// Deep copy the Params for another use case
func (p Params) Clone() Params {
	p.init()
	cpy := make(Params)
	for k, v := range p {
		cpy[k] = v
	}
	return cpy
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (p Params) Get(key string) string {
	return url.Values(p).Get(key)
}

// Set sets the key to value. It replaces any existing
// values.
func (p Params) Set(key string, value string) {
	p.init()
	url.Values(p).Set(key, value)
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (p Params) Add(key string, value string) {
	p.init()
	url.Values(p).Add(key, value)
}

// Del deletes the values associated with key.
func (p Params) Del(key string) {
	url.Values(p).Del(key)
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
func (p Params) Encode() string {
	return url.Values(p).Encode()
}

// Values returns the Params object casted to an url.Values
func (p Params) Values() url.Values {
	return url.Values(p)
}

// Merge merges the given url.Values params into the object
//
// Note that this method changes the receiver in that it adds the values from
// params into the receiver. Also, no data in the receiver are lost
func (p Params) Merge(params url.Values) Params {
	p.init()
	for k, values := range params {
		for _, v := range values {
			p.Add(k, v)
		}
	}
	return p
}

// ForTimeframe adds query params for the given timeframe
func (p *Params) ForTimeframe(timeframe Timeframe) *Params {
	p.init()
	p.Merge(timeframe.ToQuery())
	return p
}

func (p *Params) Billable(billable bool) *Params {
	p.init()
	var billableParam string
	if billable {
		billableParam = "yes"
	} else {
		billableParam = "no"
	}
	p.Set("billable", billableParam)
	return p
}

func (p *Params) OnlyBilled() *Params {
	p.init()
	p.Set("only_billed", "yes")
	return p
}

func (p *Params) OnlyUnbilled() *Params {
	p.init()
	p.Set("only_unbilled", "yes")
	return p
}

func (p *Params) IsClosed(closed bool) *Params {
	p.init()
	var isClosed string
	if closed {
		isClosed = "yes"
	} else {
		isClosed = "no"
	}
	p.Set("is_closed", isClosed)
	return p
}

func (p *Params) UpdatedSince(t time.Time) *Params {
	p.init()
	p.Set("updated_since", t.UTC().String())
	return p
}

func (p *Params) ForProject(project *Project) *Params {
	p.init()
	p.Set("project_id", string(project.Id()))
	return p
}

func (p *Params) ForUser(user *User) *Params {
	p.init()
	p.Set("user_id", string(user.Id()))
	return p
}

func (p *Params) ByClient(client *Client) *Params {
	p.init()
	p.Set("client", string(client.Id()))
	return p
}

func (p *Params) Page(page int) *Params {
	p.init()
	p.Set("page", string(page))
	return p
}

// Available status:
//   open    - sent to the client but no payment recieved
//   partial - partial payment was recorded
//   draft   - Harvest did not sent this to a client, nor recorded any payments
//   paid    - invoice paid in full
//   unpaid  - unpaid invoices
//   pastdue - past due invoices
func (p *Params) Status(status string) *Params {
	p.init()
	p.Set("status", status)
	return p
}
