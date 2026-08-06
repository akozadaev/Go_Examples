package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AMQPURL         string
	Prefetch        int
	RetryDelay      time.Duration
	MaxRetries      int
	PublishTimeout  time.Duration
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		AMQPURL:         envOrDefault("AMQP_URL", "amqp://app:app@localhost:5672/app"),
		RetryDelay:      5 * time.Second,
		PublishTimeout:  10 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}

	var err error
	if cfg.Prefetch, err = positiveInt("PREFETCH", 8); err != nil {
		return Config{}, err
	}
	if cfg.MaxRetries, err = nonNegativeInt("MAX_RETRIES", 3); err != nil {
		return Config{}, err
	}
	if cfg.RetryDelay, err = positiveDuration("RETRY_DELAY", cfg.RetryDelay); err != nil {
		return Config{}, err
	}
	if cfg.PublishTimeout, err = positiveDuration("PUBLISH_TIMEOUT", cfg.PublishTimeout); err != nil {
		return Config{}, err
	}
	if cfg.ShutdownTimeout, err = positiveDuration("SHUTDOWN_TIMEOUT", cfg.ShutdownTimeout); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func positiveInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer, got %q", key, value)
	}
	return parsed, nil
}

func nonNegativeInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer, got %q", key, value)
	}
	return parsed, nil
}

func positiveDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration, got %q", key, value)
	}
	return parsed, nil
}
