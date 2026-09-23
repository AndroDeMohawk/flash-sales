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
	consumerHandler "github.com/AndroDeMohawk/flash-sales/internal/handler/consumer"
	infraKafka "github.com/AndroDeMohawk/flash-sales/internal/infrastructure/kafka"
	infraPostgres "github.com/AndroDeMohawk/flash-sales/internal/infrastructure/postgres"
	infraRedis "github.com/AndroDeMohawk/flash-sales/internal/infrastructure/redis"
	"github.com/AndroDeMohawk/flash-sales/internal/usecase"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading config", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbPool, err := infraPostgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Error connecting to database", zap.Error(err))
	}
	defer dbPool.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer redisClient.Close()

	orderRepo := infraPostgres.NewOrderRepository(dbPool)
	statusCache := infraRedis.NewStatusRepository(redisClient)

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	kafkaConsumer := infraKafka.NewOrderConsumer(brokers, "orders.created", "flash-sales-workers")
	defer kafkaConsumer.Close()

	processUC := usecase.NewProcessOrderUseCase(orderRepo, statusCache)
	handler := consumerHandler.NewOrderConsumerHandler(kafkaConsumer, processUC, logger)

	go handler.Start(ctx)
	logger.Info("Worker Service background consumer started successfully")

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: httpMux,
	}

	go func() {
		logger.Info("Starting Worker HTTP server", zap.String("port", cfg.Port))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down Worker Service gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Error shutting down HTTP server", zap.Error(err))
	}

	logger.Info("Worker stopped gracefully")
}
