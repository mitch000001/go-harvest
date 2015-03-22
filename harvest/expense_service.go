package harvest

import (
	"fmt"
	"net/url"
)

type ExpenseService struct {
	endpoint AllEndpoint
}

func NewExpenseService(endpoint AllEndpoint) *ExpenseService {
	return &ExpenseService{endpoint: endpoint}
}

func (e *ExpenseService) All(expenses *[]*Expense, params url.Values) error {
	if len(params) == 0 || params.Get("from") == "" || params.Get("to") == "" {
		return fmt.Errorf("Bad Request: 'from' and 'to' query parameter are not optional!")
	}
	return e.endpoint.All(expenses, params)
}
