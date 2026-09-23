package domain

import (
	"context"
	"time"

	"github.com/AndroDeMohawk/flash-sales/pkg/events"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, id string) (*Order, error)
}

type StockRepository interface {
	ReserveStock(ctx context.Context, ticketID int64, quantity int) (int64, error)
	RestoreStock(ctx context.Context, ticketID int64, quantity int) error
}

type StatusRepository interface {
	SetStatus(ctx context.Context, orderID string, status string, ttl time.Duration) error
	GetStatus(ctx context.Context, orderID string) (string, error)
}
type OrderProducer interface {
	SendOrderCreated(ctx context.Context, event events.OrderCreatedEvent) error
}
