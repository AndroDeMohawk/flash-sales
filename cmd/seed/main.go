package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AndroDeMohawk/flash-sales/internal/config"
	infraPostgres "github.com/AndroDeMohawk/flash-sales/internal/infrastructure/postgres"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// 1. Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	// 2. Подключаемся к PostgreSQL
	dbPool, err := infraPostgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	// 3. Данные тестового билета
	var ticketID int64 = 1
	title := "Концерт Metallica"
	totalStock := 100
	price := 5000

	// Создаем запись в PostgreSQL (если уже есть — обновляем)
	query := `
		INSERT INTO tickets (id, title, total_stock, price)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE 
		SET total_stock = EXCLUDED.total_stock, price = EXCLUDED.price;
	`
	_, err = dbPool.Exec(ctx, query, ticketID, title, totalStock, price)
	if err != nil {
		log.Fatalf("Ошибка засева в PostgreSQL: %v", err)
	}
	fmt.Printf("[PostgreSQL] Билет ID=%d ('%s') успешно создан/обновлен!\n", ticketID, title)

	// 4. Подключаемся к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer redisClient.Close()

	// 5. Записываем остаток в Redis
	redisKey := fmt.Sprintf("ticket:%d:stock", ticketID)
	err = redisClient.Set(ctx, redisKey, totalStock, 0).Err()
	if err != nil {
		log.Fatalf("Ошибка засева в Redis: %v", err)
	}
	fmt.Printf("[Redis] Ключ '%s' успешно инициализирован со стоком: %d\n", redisKey, totalStock)

	fmt.Println("\nГотово! База и кеш забиты тестовыми данными.")
}
