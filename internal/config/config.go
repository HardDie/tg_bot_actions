package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/HardDie/tg_bot_actions/internal/logger"
)

type Config struct {
	Bot   Bot
	Cache Cache
}

func Get() Config {
	if err := godotenv.Load(); err != nil {
		if check := os.IsNotExist(err); !check {
			logger.Error.Fatalf("failed to load env vars: %s", err)
		}
	}

	cfg := Config{
		Bot:   botConfig(),
		Cache: cacheConfig(),
	}
	return cfg
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		logger.Error.Fatalf("env %q value not found", key)
	}
	return value
}
