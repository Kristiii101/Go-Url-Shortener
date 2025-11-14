package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            int
	BaseURL         string
	DatabaseURL     string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	RateLimitCreate int           // requests per minute per IP for POST /v1/links
	RateLimitWindow time.Duration // e.g., 1m
	KeyMinLen       int
	KeyMaxLen       int
	WebDir          string
}

func Load() (Config, error) {
	cfg := Config{
		Port:            intFromEnv("PORT", 8080),
		BaseURL:         strFromEnv("BASE_URL", "http://localhost:8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		ReadTimeout:     durationFromEnv("READ_TIMEOUT", 5*time.Second),
		WriteTimeout:    durationFromEnv("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:     durationFromEnv("IDLE_TIMEOUT", 60*time.Second),
		RateLimitCreate: intFromEnv("RATE_LIMIT_CREATE", 10),
		RateLimitWindow: durationFromEnv("RATE_LIMIT_WINDOW", time.Minute),
		KeyMinLen:       intFromEnv("KEY_MIN_LEN", 6),
		KeyMaxLen:       intFromEnv("KEY_MAX_LEN", 8),
		WebDir:          strFromEnv("WEB_DIR", "web"),
	}
	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("DATABASE_URL is required")
	}
	return cfg, nil
}

func strFromEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func intFromEnv(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func durationFromEnv(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
