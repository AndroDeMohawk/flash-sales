package events

import "time"

type OrderCreatedEvent struct {
	OrderID   string    `json:"order_id"`
	UserID    int64     `json:"user_id"`
	TicketID  int64     `json:"ticket_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}
