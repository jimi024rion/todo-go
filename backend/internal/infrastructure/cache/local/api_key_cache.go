package local

import (
	"context"
	"time"

	gocache "github.com/patrickmn/go-cache"

	apikeyCache "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/cache"
	apikeyentity "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/entity"
)

const (
	defaultTTL      = 5 * time.Minute
	cleanupInterval = 10 * time.Minute
)

type localAPIKeyCache struct {
	cache *gocache.Cache
}

func NewLocalAPIKeyCache() apikeyCache.APIKeyCache {
	return &localAPIKeyCache{
		cache: gocache.New(defaultTTL, cleanupInterval),
	}
}

func (c *localAPIKeyCache) Get(_ context.Context, key string) (*apikeyentity.APIKey, bool) {
	v, ok := c.cache.Get(key)
	if !ok {
		return nil, false
	}
	return v.(*apikeyentity.APIKey), true
}

func (c *localAPIKeyCache) Set(_ context.Context, key string, apiKey *apikeyentity.APIKey) {
	c.cache.Set(key, apiKey, gocache.DefaultExpiration)
}

func (c *localAPIKeyCache) Delete(_ context.Context, key string) {
	c.cache.Delete(key)
}
