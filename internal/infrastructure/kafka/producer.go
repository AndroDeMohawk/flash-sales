package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AndroDeMohawk/flash-sales/pkg/events"
	"github.com/segmentio/kafka-go"
)

type OrderProducer struct {
	writer *kafka.Writer
}

func NewOrderProducer(brokers []string, topic string) *OrderProducer {
	return &OrderProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.Hash{},
			Topic:        topic,
			RequiredAcks: kafka.RequireOne,
			Async:        true,
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *OrderProducer) SendOrderCreated(ctx context.Context, event events.OrderCreatedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal order created event: %w", err)
	}
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", event.TicketID)),
		Value: payload,
		Time:  time.Now(),
	}
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send order created event: %w", err)
	}
	return nil
}

func (p *OrderProducer) Close() error {
	return p.writer.Close()
}
