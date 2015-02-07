package harvest

import "time"

type DayEntry struct {
	Hours     float64   `json:"hours"`
	ID        int       `json:"id"`
	Notes     string    `json:"notes"`
	ProjectId int       `json:"project-id"`
	SpentAt   time.Time `json:"spent-at"`
	TaskId    int       `json:"task-id"`
	UserId    int       `json:"user-id"`
	// was this record invoiced, or marked as invoiced
	IsBilled bool `json:"is-billed"`
	// was this record approved or not (for accounts with approval feature)
	IsClosed  bool      `json:"is-closed"`
	UpdatedAt time.Time `json:"updated-at"`
	CreatedAt time.Time `json:"created-at"`
}
