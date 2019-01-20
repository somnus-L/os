package cache
//----------------------------------------------*
// 基于 go-redis
//----------------------------------------------*


import (
	"github.com/go-redis/redis"
	"time"
	"compress/gzip"
	"context"
	"reflect"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{Client: client}
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	o, ok := value.(BinaryMarshal)
	if !ok {
		o = &DefaultMarshal{value, gzip.NoCompression}
	}
	data, err := o.Marshal(ctx)
	if err != nil {
		return err
	}
	return c.Client.Set(key, data, ttl).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) (interface{}, error) {
	stringValue, err := c.Client.Get(key).Result()
	if err == redis.Nil {
		return nil, ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	o, ok := value.(BinaryUnmarshal)
	if !ok {
		o = &DefaultMarshal{value, gzip.NoCompression}
	}
	err = o.Unmarshal(ctx, []byte(stringValue))
	if err != nil {
		return nil, err
	}
	if value == nil {
		return value, nil
	}
	return reflect.ValueOf(value).Elem().Interface(), nil
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.Client.Del(key).Err()
}

