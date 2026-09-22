package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type OrderConsumer struct {
	reader *kafka.Reader
}

func NewOrderConsumer(brokers []string, topic, groupID string) *OrderConsumer {
	return &OrderConsumer{reader: kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10,
		MaxBytes:       10e6,
		MaxWait:        1 * time.Second,
		CommitInterval: 0,
		StartOffset:    kafka.FirstOffset,
	}),
	}
}

func (c *OrderConsumer) FetchMessage(ctx context.Context) (error, kafka.Message) {
	msg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return fmt.Errorf("FetchMessage: %w", err), kafka.Message{}
	}
	return nil, msg

}

func (c *OrderConsumer) CommitMessage(ctx context.Context, msg kafka.Message) error {
	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		return fmt.Errorf("Error to commit message %w", err)
	}
	return nil
}

func (c *OrderConsumer) Close() error {
	return c.reader.Close()
}
