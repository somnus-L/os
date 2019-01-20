package cache

import (
	"testing"
	"context"
	"time"
	"strconv"
)

func TestMemoryCache(t *testing.T) {
	var (
		c   Cache = NewMemoryCache(10)
		ctx       = context.Background()
		err error
	)

	str1 := "abc"
	str2 := ""
	err = c.Set(ctx, "test", str1, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	str2Itf, err := c.Get(ctx, "test", &str2)
	if err != nil {
		t.Fatal(err)
	}
	if str2Itf.(string) != str1 {
		t.Fatal("get error")
	}
	err = c.Del(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}
	str2Itf, err = c.Get(ctx, "test", &str2)
	if err != ErrNotExist {
		t.Fatal("del error")
	}
}

func TestMemoryCache_GC(t *testing.T) {
	var (
		c   Cache = NewMemoryCache(10)
		ctx       = context.Background()
		err error
	)
	for i := 0; i < 100; i++ {
		err = c.Set(ctx, strconv.Itoa(i), i, time.Second)
		if err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < 100; i++ {
		_, err = c.Get(ctx, strconv.Itoa(i), nil)
		if err != ErrNotExist {
			t.Fatal("gc error")
		}
	}
}

func TestMemoryCache_Max(t *testing.T) {
	var (
		max = 5
		c   = NewMemoryCache(max)
		ctx = context.Background()
		err error
	)
	for i := 0; i < max*10; i++ {
		err = c.Set(ctx, strconv.Itoa(i), i, time.Hour)
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(c.m) > max || c.l.Len() > max {
		t.Fatal("max error")
	}
	t.Log(len(c.m), c.l.Len())
	for _, e := range c.m {
		n := e.Value.(*node)
		t.Log(n.key, n.value)
	}
}
