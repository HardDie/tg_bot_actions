package repository

import (
	"github.com/jellydator/ttlcache/v3"

	"github.com/HardDie/tg_bot_actions/internal/config"
	"github.com/HardDie/tg_bot_actions/internal/models"
	"github.com/HardDie/tg_bot_actions/internal/utils"
)

type CacheRepository struct {
	cache *ttlcache.Cache[string, models.Cache]
}

func NewCacheRepository(cfg config.Cache) *CacheRepository {
	return &CacheRepository{
		cache: ttlcache.New[string, models.Cache](
			ttlcache.WithTTL[string, models.Cache](cfg.Period),
			ttlcache.WithDisableTouchOnHit[string, models.Cache](),
		),
	}
}

func (r CacheRepository) Get(key string) *models.Cache {
	item := r.cache.Get(key)
	if item == nil {
		return nil
	}
	return utils.Allocate(item.Value())
}

func (r CacheRepository) Set(key string, value models.Cache) {
	r.cache.Set(key, value, ttlcache.DefaultTTL)
}
