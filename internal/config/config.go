package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string

	AppPort string

	HELA_TESTNET_RPC string
	WSSURL           string
	ContractAddress  string

	HTTPReadHeaderTimeout time.Duration
	HTTPReadTimeout       time.Duration
	HTTPWriteTimeout      time.Duration
	HTTPIdleTimeout       time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),

		AppPort: getEnv("APP_PORT", "8080"),

		HELA_TESTNET_RPC: getEnv("HELA_TESTNET_RPC", ""),
		WSSURL:           getEnv("WSS_URL", ""),
		ContractAddress:  getEnv("CONTRACT_ADDRESS", ""),

		HTTPReadHeaderTimeout: getDurationSeconds("HTTP_READ_HEADER_TIMEOUT_SECONDS", 5),
		HTTPReadTimeout:       getDurationSeconds("HTTP_READ_TIMEOUT_SECONDS", 10),
		HTTPWriteTimeout:      getDurationSeconds("HTTP_WRITE_TIMEOUT_SECONDS", 10),
		HTTPIdleTimeout:       getDurationSeconds("HTTP_IDLE_TIMEOUT_SECONDS", 60),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.HELA_TESTNET_RPC == "" {
		return nil, fmt.Errorf("HELA_TESTNET_RPC is required")
	}

	if cfg.WSSURL == "" {
		return nil, fmt.Errorf("WSS_URL is required")
	}

	if !common.IsHexAddress(cfg.ContractAddress) {
		return nil, fmt.Errorf("CONTRACT_ADDRESS is invalid")
	}

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getDurationSeconds(key string, fallback int) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(fallback) * time.Second
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return time.Duration(fallback) * time.Second
	}

	return time.Duration(seconds) * time.Second
}
