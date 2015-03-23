package harvest

import "fmt"

func (i *InvoiceService) PublicURL(invoice *Invoice) string {
	url := i.endpoint.URL()
	return fmt.Sprintf("%s/client/invoices/%s", url.String(), invoice.ClientKey)
}
