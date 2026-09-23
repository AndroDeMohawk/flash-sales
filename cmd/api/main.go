package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/AndroDeMohawk/flash-sales/internal/config"
	httpHandler "github.com/AndroDeMohawk/flash-sales/internal/handler/http"
	infraKafka "github.com/AndroDeMohawk/flash-sales/internal/infrastructure/kafka"
	infraRedis "github.com/AndroDeMohawk/flash-sales/internal/infrastructure/redis"
	"github.com/AndroDeMohawk/flash-sales/internal/usecase"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	//Логгер
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	//Конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading config", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	//Редис
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer redisClient.Close()

	//Репо
	stockRepo := infraRedis.NewStockRepository(redisClient)
	statusCache := infraRedis.NewStatusRepository(redisClient)

	//Кафка
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	kafkaProducer := infraKafka.NewOrderProducer(brokers, "orders.created")
	defer kafkaProducer.Close()

	//Сборка
	orderUC := usecase.NewOrderUseCase(stockRepo, statusCache, kafkaProducer)
	orderHandler := httpHandler.NewOrderHandler(orderUC, logger)
	mainHandler := httpHandler.NewHandler(orderHandler, logger)

	//Сервер
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mainHandler.InitRoutes(),
	}

	//Старт
	go func() {
		logger.Info("Starting API Service HTTP server", zap.String("port", cfg.Port))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down API Service gracefully...")

	//Шатдаун
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Error shutting down HTTP server", zap.Error(err))
	}

	logger.Info("API Service stopped gracefully")
}
