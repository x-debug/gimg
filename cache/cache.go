package cache

import (
	"errors"
	"gimg/config"
)

var CacheMiss = errors.New("Cache miss")
var CacheNotFound = errors.New("Cache not found")

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, values []byte) error
}

func NewCache(cfg *config.CacheConf) (Cache, error) {
	if cfg.Type == "memcached" {
		c := &memcachedCache{}
		err := c.initialize(cfg.Brockers)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	return nil, CacheNotFound
}
