package cache

import (
	"gimg/config"
	lg "gimg/logger"
	"time"

	"github.com/allegro/bigcache/v3"
)

type memoryCache struct {
	client *bigcache.BigCache
	logger lg.Logger
}

func (mc *memoryCache) initialize(logger lg.Logger, cfg *config.CacheConf) error {
	mc.logger = logger
	logger.Info("Initialize memory cache", lg.Int64("Lifetime", cfg.LifeTime))
	cli, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Duration(cfg.LifeTime) * time.Minute))
	if err != nil {
		return err
	}
	mc.client = cli
	return nil
}

func (mc *memoryCache) Set(key string, values []byte) error {
	return mc.client.Set(key, values)
}

func (mc *memoryCache) Get(key string) ([]byte, error) {
	item, err := mc.client.Get(key)
	if err == bigcache.ErrEntryNotFound {
		return nil, CacheMiss
	}

	return item, err
}
