package cache

import (
	"testing"
	"compress/gzip"
	"context"
	"time"
	"gitlab.followme.com/CopyTradingGo/guard-sdk/builder"
)

func TestRedisCache(t *testing.T) {
	b := builder.RedisBuilder{
		Addr:            "127.0.0.1:6379",
		Password:        "",
		MaxRetries:      3,
		MinRetryBackoff: 10 * time.Millisecond,
		MaxRetryBackoff: time.Second,
		PoolSize:        20,
		MinIdleConns:    5,
	}
	rc, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}

	var (
		c   Cache = &RedisCache{Client: rc}
		ctx       = context.Background()
	)

	// default
	str1 := "abc"
	str2 := ""
	err = c.Set(ctx, "test", str1, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Get(ctx, "test", &str2)
	if err != nil {
		t.Fatal(err)
	}
	if str2 != str1 {
		t.Fatal("get error")
	}
	err = c.Del(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Get(ctx, "test", &str2)
	if err != ErrNotExist {
		t.Fatal("del error")
	}

	// object
	str1 = "cba"
	o1 := &DefaultMarshal{&str1, gzip.BestCompression}
	o2 := &DefaultMarshal{&str2, gzip.BestCompression}
	err = c.Set(ctx, "test", o1, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Get(ctx, "test", o2)
	if err != nil {
		t.Fatal(err)
	}
	if str1 != str2 {
		t.Fatal("get error")
	}
	err = c.Del(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Get(ctx, "test", o2)
	if err != ErrNotExist {
		t.Fatal("del error")
	}
}

