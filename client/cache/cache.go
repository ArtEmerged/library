package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	HashSet(ctx context.Context, key, field string, value interface{}) error
	HashGet(ctx context.Context, key, field string) (interface{}, error)
	HashDel(ctx context.Context, key, field string) error
	HGetAll(ctx context.Context, key string) ([]interface{}, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Ping(ctx context.Context) error
	Close() error
}
