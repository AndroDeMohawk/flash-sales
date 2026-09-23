package postgres

import (
	"context"
	"fmt"

	"github.com/AndroDeMohawk/flash-sales/internal/domain"
	db "github.com/AndroDeMohawk/flash-sales/internal/infrastructure/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	_, err := r.queries.CreateOrder(ctx, db.CreateOrderParams{
		ID:       order.ID,
		UserID:   order.UserID,
		TicketID: order.TicketID,
		Quantity: int32(order.Quantity),
		Status:   order.Status,
		CreatedAt: pgtype.Timestamptz{
			Time:  order.CreatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("CreateOrder: %w", err)
	}
	return nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderId string) (*domain.Order, error) {
	order, err := r.queries.GetOrderByID(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("GetOrderByID: %w", err)
	}
	return &domain.Order{
		ID:        order.ID,
		UserID:    order.UserID,
		TicketID:  order.TicketID,
		Quantity:  int(order.Quantity),
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Time,
	}, nil
}
