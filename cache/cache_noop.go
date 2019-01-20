package cache

import (
	"context"
	"time"
)

type NoopCache struct{}

func NewNoopCache() Cache {
	return &NoopCache{}
}

func (c *NoopCache) Get(ctx context.Context, key string, value interface{}) (interface{}, error) {
	return nil, ErrNotExist
}

func (c *NoopCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (c *NoopCache) Del(ctx context.Context, key string) error {
	return nil
}
