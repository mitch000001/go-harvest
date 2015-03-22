package harvest

import "time"

type Expense struct {
	ExpenseCategoryId float64   `json:"expense-category-id"`
	ID                int       `json:"id"`
	Notes             string    `json:"notes"`
	ProjectId         int       `json:"project-id"`
	UserId            int       `json:"user-id"`
	SpentAt           time.Time `json:"spent-at"`
	TotalCost         float64   `json:"total-cost"`
	Units             float64   `json:"units"`
	TaskId            int       `json:"task-id"`
	// was this record invoiced, or marked as invoiced
	IsBilled bool `json:"is-billed"`
	// was this record approved or not (for accounts with approval feature)
	IsClosed  bool      `json:"is-closed"`
	UpdatedAt time.Time `json:"updated-at"`
	CreatedAt time.Time `json:"created-at"`
}
