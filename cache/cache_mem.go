package cache

import (
	"sync"
	"context"
	"time"
	"container/list"
	"github.com/opentracing/opentracing-go"
)

// 1.maxElementSize
// 2.ttl
// 3.lru remove

type node struct {
	key        string
	value      interface{}
	expireTime time.Time
}

type MemoryCache struct {
	sync.Mutex
	l              *list.List               // element type => node
	m              map[string]*list.Element // element type => node
	maxElementSize int
}

// maxElementSize<=0 代表没有最大数量限制
func NewMemoryCache(maxElementSize int) *MemoryCache {
	mc := &MemoryCache{
		m:              make(map[string]*list.Element),
		l:              list.New(),
		maxElementSize: maxElementSize,
	}
	mc.runGC()
	return mc
}

func (c *MemoryCache) runGC() {
	time.AfterFunc(time.Second*5, func() {
		c.runGC()
		c.GC()
	})
}

const GcSampleCount = 20
const Gc_Sample_Percentage = 0.25

func (c *MemoryCache) GC() {
	c.Lock()
	defer c.Unlock()

	for {
		var count float64 = 1
		var gcCount float64 = 0.0
		gcStart := time.Now()
		for _, e := range c.m {
			count++
			n := e.Value.(*node)
			if n.expireTime.Before(gcStart) {
				c.l.Remove(e)
				delete(c.m, n.key)
				gcCount++
			}
			if count > GcSampleCount {
				break
			}
		}
		if gcCount/count < Gc_Sample_Percentage {
			break
		}
	}
}

func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "MemoryCache.Set")
	defer span.Finish()

	var e *list.Element
	var ok bool

	c.Lock()
	defer c.Unlock()

	e, ok = c.m[key]
	if ok {
		e.Value.(*node).expireTime = time.Now().Add(ttl)
	} else {
		e = c.l.PushBack(&node{key, value, time.Now().Add(ttl)})
		c.m[key] = e
	}

	if c.maxElementSize > 0 && c.l.Len() > c.maxElementSize {
		e := c.l.Front()
		n := c.l.Remove(e).(*node)
		delete(c.m, n.key)
	}
	return nil
}

func (c *MemoryCache) Get(ctx context.Context, key string, value interface{}) (interface{}, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "MemoryCache.Get")
	defer span.Finish()

	var e *list.Element
	var ok bool

	c.Lock()
	defer c.Unlock()

	e, ok = c.m[key]
	if !ok {
		return nil, ErrNotExist
	}

	n := e.Value.(*node)
	if n.expireTime.Before(time.Now()) {
		c.l.Remove(e)
		delete(c.m, n.key)
		return nil, ErrNotExist
	}
	c.l.MoveToBack(e)
	return n.value, nil
}

func (c *MemoryCache) Del(ctx context.Context, key string) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "MemoryCache.Del")
	defer span.Finish()

	c.Lock()
	defer c.Unlock()

	e, ok := c.m[key]
	if !ok {
		return nil
	}
	c.l.Remove(e)
	delete(c.m, key)
	return nil
}
