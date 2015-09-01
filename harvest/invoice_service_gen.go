// DO NOT EDIT!
// This file is generated by the api generator.

// +build !feature

package harvest

import (
	"net/url"
)

type InvoiceService struct {
	endpoint	CrudEndpoint
	provider	CrudEndpointProvider
}

func NewInvoiceService(endpoint CrudEndpoint, provider CrudEndpointProvider) *InvoiceService {
	service := InvoiceService{
		endpoint:	endpoint,
		provider:	provider,
	}
	return &service
}

func (s *InvoiceService) All(invoices *[]*Invoice, params url.Values) error {
	return s.endpoint.All(invoices, params)
}

func (s *InvoiceService) Find(id int, invoice *Invoice, params url.Values) error {
	return s.endpoint.Find(id, invoice, params)
}

func (s *InvoiceService) Create(invoice *Invoice) error {
	return s.endpoint.Create(invoice)
}

func (s *InvoiceService) Update(invoice *Invoice) error {
	return s.endpoint.Update(invoice)
}

func (s *InvoiceService) Delete(invoice *Invoice) error {
	return s.endpoint.Delete(invoice)
}
