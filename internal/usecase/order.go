package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/AndroDeMohawk/flash-sales/internal/domain"
	"github.com/AndroDeMohawk/flash-sales/internal/infrastructure/redis"
	"github.com/AndroDeMohawk/flash-sales/pkg/events"
)

type OrderUseCase struct {
	stockRepo   domain.StockRepository
	statusCache domain.StatusRepository
	producer    domain.OrderProducer
}

func NewOrderUseCase(
	stockRepo domain.StockRepository,
	statusCache domain.StatusRepository,
	producer domain.OrderProducer,
) *OrderUseCase {
	return &OrderUseCase{
		stockRepo:   stockRepo,
		statusCache: statusCache,
		producer:    producer,
	}
}
func (uc *OrderUseCase) CreateOrder(ctx context.Context, userId, ticketId int64, quantity int) (string, error) {
	_, err := uc.stockRepo.ReserveStock(ctx, ticketId, quantity)
	if err != nil {
		if errors.Is(err, redis.ErrSoldOut) {
			return "", redis.ErrSoldOut
		}
		return "", fmt.Errorf("failed to reserve stock: %w", err)
	}
	orderId := uuid.New().String()

	_ = uc.statusCache.SetStatus(ctx, orderId, "PENDING", 24*time.Hour)
	event := events.OrderCreatedEvent{
		OrderID:   orderId,
		UserID:    userId,
		TicketID:  ticketId,
		Quantity:  quantity,
		CreatedAt: time.Now(),
	}
	if err := uc.producer.SendOrderCreated(ctx, event); err != nil {
		return "", fmt.Errorf("failed to send order created: %w", err)
	}
	return orderId, nil
}
func (uc *OrderUseCase) GetOrderStatus(ctx context.Context, orderID string) (string, error) {
	status, err := uc.statusCache.GetStatus(ctx, orderID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch order status: %w", err)
	}
	return status, nil
}
