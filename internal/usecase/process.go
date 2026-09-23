package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/AndroDeMohawk/flash-sales/internal/domain"
	"github.com/AndroDeMohawk/flash-sales/pkg/events"
)

type ProcessOrderUseCase struct {
	orderRepo   domain.OrderRepository
	statusCache domain.StatusRepository
}

func NewProcessOrderUseCase(orderRepo domain.OrderRepository, statusCache domain.StatusRepository) *ProcessOrderUseCase {
	return &ProcessOrderUseCase{
		orderRepo:   orderRepo,
		statusCache: statusCache,
	}
}

func (uc *ProcessOrderUseCase) Process(ctx context.Context, event events.OrderCreatedEvent) error {
	order := &domain.Order{
		ID:        event.OrderID,
		UserID:    event.UserID,
		TicketID:  event.TicketID,
		Quantity:  event.Quantity,
		Status:    "CONFIRMED",
		CreatedAt: event.CreatedAt,
	}
	if err := uc.orderRepo.CreateOrder(ctx, order); err != nil {
		return fmt.Errorf("error to create order: %w", err)
	}
	if err := uc.statusCache.SetStatus(ctx, event.OrderID, "CONFIRMED", 24*time.Hour); err != nil {
		return fmt.Errorf("failed to update status redis: %w", err)
	}
	return nil

}
