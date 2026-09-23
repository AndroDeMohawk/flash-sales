package domain

import "time"

type Order struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	TicketID  int64     `json:"ticket_id"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
