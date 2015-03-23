package harvest

import "time"

//go:generate go run ../cmd/api_gen/api_gen.go -type=Invoice -c -fields CrudEndpointProvider

type Invoice struct {
	ID        int       `json:"id"`
	Amount    float64   `json:"amount"`
	DueAmount float64   `json:"due-amount"`
	DueAt     ShortDate `json:"due-at"`
	// human representation for due at
	DueAtHumanFormat string `json:"due-at-human-format"`
	// invoiced period, present for generated invoices
	PeriodEnd   ShortDate `json:"period-end"`
	PeriodStart ShortDate `json:"period-start"`
	ClientId    int       `json:"client-id"`
	Subject     string    `json:"subject"`
	// see
	Currency      string    `json:"currency"`
	IssuedAt      ShortDate `json:"issued-at"`
	CreatedById   int       `json:"created-by-id"`
	Notes         string    `json:"notes"`
	Number        string    `json:"number"`
	PurchaseOrder string    `json:"purchase-order"`
	ClientKey     string    `json:"client-key"`
	// See invoice messages and invoice payments for manipulating the state attribute.  Direct assigment will be ignored. Options are open, draft, partial, paid and closed
	State string `json:"state"`
	// applied tax percentage, blank if not taxed
	Tax float64 `json:"tax"`
	// applied tax 2 percentage, blank if not taxed
	Tax2 float64 `json:"tax2"`
	// the first tax amount
	TaxAmount float64 `json:"tax-amount"`
	// the second tax amount
	TaxAmount2 float64 `json:"tax-amount2"`
	// discount
	DiscountAmount float64 `json:"discount-amount"`
	Discount       float64 `json:"discount"`
	// is it recurring?
	RecurringInvoiceId int `json:"recurring-invoice-id"`
	// was this converted from an estimate?
	EstimateId int `json:"estimate-id"`
	// a retainer_id will only be present if the invoice funds a retainer
	RetainerId   int       `json:"retainer-id"`
	UpdatedAt    time.Time `json:"updated-at"`
	CreatedAt    time.Time `json:"created-at"`
	CsvLineItems string    `json:"csv-line-items"`
	/* allowed values:
	   free_form:  creates a free form invoice (non-generated).
	               Content is added via CSV.
	   project:    gathers content from Harvest grouping by projects
	   task:       gathers content from Harvest grouping by task
	   people:     gathers content from Harvest grouping by people
	   detailed:   includes detailed notes */
	Kind string `json:"kind"`
	// comma separated project ids to gather data from, useless on free_form invoices
	ProjectsToInvoice string `json:"projects-to-invoice"`
	// import hours useless on free_form invoices
	ImportHours string `json:"import-hours"`
	// import expenses useless on free_form invoices
	ImportExpenses string `json:"import-expenses"`
	// invoiced period for expenses, present for generated invoices
	ExpensePeriodEnd   ShortDate `json:"expense-period-end"`
	ExpensePeriodStart ShortDate `json:"expense-period-start"`
}

func (i *Invoice) Id() int {
	return i.ID
}

func (i *Invoice) SetId(id int) {
	i.ID = id
}

func (i *Invoice) Type() string {
	return "Invoice"
}
