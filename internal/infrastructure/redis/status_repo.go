package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type StatusRepository struct {
	client *redis.Client
}

func NewStatusRepository(client *redis.Client) *StatusRepository {
	return &StatusRepository{client: client}
}

func (r *StatusRepository) SetStatus(ctx context.Context, OrderId, status string, ttl time.Duration) error {
	key := fmt.Sprintf("Order:%s:status:", OrderId)
	if err := r.client.Set(ctx, key, status, ttl).Err(); err != nil {
		return fmt.Errorf("SetStatus: %w", err)
	}
	return nil
}
func (r *StatusRepository) GetStatus(ctx context.Context, OrderId string) (string, error) {
	key := fmt.Sprintf("Order:%s:status:", OrderId)
	status, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("status foe order %s not exists", OrderId)
		}
		return "", fmt.Errorf("GetStatus: %w", err)
	}
	return status, nil
}
