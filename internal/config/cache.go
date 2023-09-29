package config

import (
	"time"

	"github.com/HardDie/tg_bot_actions/internal/utils"
)

type Cache struct {
	Period time.Duration
}

func cacheConfig() Cache {
	return Cache{
		Period: utils.Must(time.ParseDuration(getEnv("CACHE_PERIOD"))),
	}
}
