package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/AndroDeMohawk/flash-sales/internal/infrastructure/kafka"
	"github.com/AndroDeMohawk/flash-sales/internal/usecase"
	"github.com/AndroDeMohawk/flash-sales/pkg/events"
	"go.uber.org/zap"
)

type OrderConsumerHandler struct {
	consumer  *kafka.OrderConsumer
	ProcessUC *usecase.ProcessOrderUseCase
	logger    *zap.Logger
}

func NewOrderConsumerHandler(
	consumer *kafka.OrderConsumer,
	processUC *usecase.ProcessOrderUseCase,
	logger *zap.Logger,
) *OrderConsumerHandler {
	return &OrderConsumerHandler{
		consumer:  consumer,
		ProcessUC: processUC,
		logger:    logger,
	}
}
func (h *OrderConsumerHandler) Start(ctx context.Context) {
	h.logger.Info("Start Kafka Order Consumer loop")
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("Kafka Order Consumer loop stopped")
			return
		default:
			msg, err := h.consumer.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				h.logger.Error("Error to fetching message from Kafka", zap.Error(err))
				continue
			}
			var event events.OrderCreatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				h.logger.Error("Error to unmarshal event", zap.Error(err))
				_ = h.consumer.CommitMessage(ctx, msg)
				continue
			}
			if err := h.ProcessUC.Process(ctx, event); err != nil {
				h.logger.Error("Error to process event", zap.String("order_id", event.OrderID), zap.Error(err))
				continue
			}
			if err := h.consumer.CommitMessage(ctx, msg); err != nil {
				h.logger.Error("Error to commit event", zap.Error(err))
			} else {
				h.logger.Info("Commit event", zap.String("order_id", event.OrderID))
			}
		}

	}
}
