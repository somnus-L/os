package cache

import (
	"errors"
	"context"
	"github.com/opentracing/opentracing-go"
	"time"
)

var ErrNotExist = errors.New("not exist")

//Set
// var a struct A
// cache.Set(ctx, "example", a)
// var b struct A

//Get
//b Must be assigned to ensure deserialization
//b = cache.Get(ctx, "example", b)
//Get requires an additional value parameter to be passed in to be compatible with the binary serialized cache


type Cache interface {
	Get(ctx context.Context, key string, value interface{}) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type BinaryMarshal interface {
	Marshal(ctx context.Context) (data []byte, err error)
}

type BinaryUnmarshal interface {
	Unmarshal(ctx context.Context, data []byte) error
}

func DeepClone(ctx context.Context, a BinaryMarshal, b BinaryUnmarshal) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "DeepClone")
	defer span.Finish()

	data, err := a.Marshal(ctx)
	if err != nil {
		return err
	}
	return b.Unmarshal(ctx, data)
}
