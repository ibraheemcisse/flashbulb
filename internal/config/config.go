package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL       string
	HTTPAddr          string
	DBMaxConns        int32
	DBMinConns        int32
	DBMaxConnLifetime time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		HTTPAddr:          getEnv("HTTP_ADDR", ":8080"),
		DBMaxConns:        10,
		DBMinConns:        2,
		DBMaxConnLifetime: time.Hour,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	maxConns, err := getInt32Env("DB_MAX_CONNS", cfg.DBMaxConns)
	if err != nil {
		return Config{}, err
	}
	cfg.DBMaxConns = maxConns

	minConns, err := getInt32Env("DB_MIN_CONNS", cfg.DBMinConns)
	if err != nil {
		return Config{}, err
	}
	cfg.DBMinConns = minConns

	lifetime, err := time.ParseDuration(
		getEnv("DB_MAX_CONN_LIFETIME", cfg.DBMaxConnLifetime.String()),
	)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_MAX_CONN_LIFETIME: %w", err)
	}
	cfg.DBMaxConnLifetime = lifetime

	if cfg.HTTPAddr == "" {
		return Config{}, fmt.Errorf("HTTP_ADDR cannot be empty")
	}

	if cfg.DBMaxConns < 1 {
		return Config{}, fmt.Errorf("DB_MAX_CONNS must be at least 1")
	}

	if cfg.DBMinConns < 0 || cfg.DBMinConns > cfg.DBMaxConns {
		return Config{}, fmt.Errorf("DB_MIN_CONNS must be between 0 and DB_MAX_CONNS")
	}

	if cfg.DBMaxConnLifetime <= 0 {
		return Config{}, fmt.Errorf("DB_MAX_CONN_LIFETIME must be positive")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getInt32Env(key string, fallback int32) (int32, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return int32(parsed), nil
}
