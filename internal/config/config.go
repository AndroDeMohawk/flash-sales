package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Version      string
	Port         string
	ENV          string
	DatabaseURL  string
	RedisAddr    string
	KafkaBrokers string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	port := getEnv("PORT", "8888")
	env := getEnv("ENV", "dev")
	version := getEnv("VERSION", "1.0.0")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	return &Config{
		Version:      version,
		Port:         port,
		ENV:          env,
		DatabaseURL:  databaseURL,
		RedisAddr:    redisAddr,
		KafkaBrokers: kafkaBrokers,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
