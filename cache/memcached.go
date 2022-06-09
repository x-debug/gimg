package cache

import (
	"strconv"

	"github.com/x-debug/gimg/config"
	"github.com/x-debug/gimg/logger"

	"github.com/bradfitz/gomemcache/memcache"
)

type memcachedCache struct {
	client *memcache.Client
	logger logger.Logger
}

func (mc *memcachedCache) initialize(logger logger.Logger, cfg *config.CacheConf) error {
	mc.logger = logger
	addrPairs := make([]string, 0)
	for _, brocker := range cfg.Brockers {
		addrPairs = append(addrPairs, brocker.Addr+":"+strconv.Itoa(brocker.Port))
	}
	mc.client = memcache.New(addrPairs...)
	return nil
}

func (mc *memcachedCache) Set(key string, values []byte) error {
	return mc.client.Set(&memcache.Item{Key: key, Value: values, Flags: 0, Expiration: 0}) //no expiration time
}

func (mc *memcachedCache) Get(key string) ([]byte, error) {
	item, err := mc.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, CacheMiss
		}

		return nil, err
	}
	return item.Value, nil
}
