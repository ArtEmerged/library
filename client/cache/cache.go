package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, in interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	HSet(ctx context.Context, key, field string, value interface{}) error
	HGet(ctx context.Context, key, field string, in interface{}) error
	HDel(ctx context.Context, key string, field ...string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Ping(ctx context.Context) error
	Close() error
}
