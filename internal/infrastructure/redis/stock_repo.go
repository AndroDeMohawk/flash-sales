package redis

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var reserveStockLua string

var (
	ErrTicketNotFound = errors.New("ticket stock not found in cache")
	ErrSoldOut        = errors.New("ticket is sold out")
)

type StockRepository struct {
	client *redis.Client
	script *redis.Script
}

func NewStockRepository(client *redis.Client) *StockRepository {
	return &StockRepository{
		client: client,
		script: redis.NewScript(reserveStockLua),
	}
}

// ReserveStock атомарно списывает билет из Redis
func (r *StockRepository) ReserveStock(ctx context.Context, ticketID int64, quantity int) (int64, error) {
	key := fmt.Sprintf("ticket:%d:stock", ticketID)

	// Выполняем Lua-скрипт в Redis
	res, err := r.script.Run(ctx, r.client, []string{key}, quantity).Int64()
	if err != nil {
		return 0, fmt.Errorf("redis lua execution error: %w", err)
	}

	switch res {
	case -1:
		return 0, ErrTicketNotFound
	case -2:
		return 0, ErrSoldOut
	default:
		return res, nil // Возвращает новый остаток
	}
}
func (r *StockRepository) RestoreStock(ctx context.Context, ticketID int64, quantity int) error {
	key := fmt.Sprintf("ticket:%d:stock", ticketID)

	if err := r.client.IncrBy(ctx, key, int64(quantity)).Err(); err != nil {
		return fmt.Errorf("failed to restore stock in redis: %w", err)
	}

	return nil
}
