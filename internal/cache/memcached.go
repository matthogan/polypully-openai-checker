package cache

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/codejago/polypully-openai-checker/internal/config"
	"log/slog"
)

var mc *memcache.Client

type Cache struct {
	config        *config.Cache
	keySerializer keySerializer
}

func NewCache(config *config.Config) *Cache {
	if config.Cache.Enabled {
		mc = memcache.New(config.Cache.Host)
	}
	ks := newKeySerializer(config)
	return &Cache{
		config:        &config.Cache,
		keySerializer: ks,
	}
}

func (c *Cache) KeyInit(value interface{}) (string, error) {
	if c.config.Enabled {
		key, err := json.Marshal(value)
		if err != nil {
			slog.Warn("unable to marshall a request for the cache key", "error", err)
			return "", err
		}
		hash, err := c.keySerializer.serialise(string(key))
		if err != nil {
			return "", fmt.Errorf("failed to serialise key: %w", err)
		}
		return hash, nil
	}
	return "", fmt.Errorf("cache disabled")
}

func (c *Cache) Set(key string, value string) error {
	if !c.config.Enabled {
		return nil
	}
	return mc.Set(&memcache.Item{Key: key, Value: []byte(value)})
}

func (c *Cache) Get(key string) ([]byte, error) {
	if !c.config.Enabled {
		return nil, nil
	}
	item, err := mc.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}
