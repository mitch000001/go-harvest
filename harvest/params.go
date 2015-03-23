package harvest

import (
	"net/url"
	"time"
)

type Params url.Values

func (p *Params) init() {
	if *p == nil {
		*p = make(Params)
	}
}

// Deep copy the Params for another use case
func (p *Params) Clone() *Params {
	cpy := make(Params)
	for k, v := range *p {
		cpy[k] = v
	}
	return &cpy
}

func (p *Params) Add(key string, value string) {
	url.Values(*p).Add(key, value)
}

func (p *Params) Set(key string, value string) {
	url.Values(*p).Set(key, value)
}

func (p *Params) Get(key string) string {
	return url.Values(*p).Get(key)
}

func (p *Params) Encode() string {
	return url.Values(*p).Encode()
}

func (p *Params) Values() url.Values {
	return url.Values(*p)
}

func (p *Params) ForTimeframe(timeframe Timeframe) *Params {
	p.init()
	p.Set("from", timeframe.StartDate.Format("20060102"))
	p.Set("to", timeframe.EndDate.Format("20060102"))
	return p
}

func (p *Params) OnlyBillable(billed bool) *Params {
	p.init()
	var billable string
	if billed {
		billable = "yes"
	} else {
		billable = "no"
	}
	p.Set("billable", billable)
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
