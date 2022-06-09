package cache

import (
	"errors"

	"github.com/x-debug/gimg/config"
	"github.com/x-debug/gimg/logger"
)

var CacheMiss = errors.New("Cache miss")
var CacheNotFound = errors.New("Cache not found")

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, values []byte) error
}

func NewCache(cfg *config.CacheConf, logger logger.Logger) (Cache, error) {
	var err error
	var cache Cache

	if cfg.Type == "memcached" {
		c := &memcachedCache{}
		err = c.initialize(logger, cfg)
		logger.Info("Memcached initialized")
		cache = c
	} else if cfg.Type == "memory" {
		c := &memoryCache{}
		err = c.initialize(logger, cfg)
		logger.Info("Memory initialized")
		cache = c
	} else {
		return nil, CacheNotFound
	}
	if err != nil {
		return nil, err
	}

	return cache, nil
}
