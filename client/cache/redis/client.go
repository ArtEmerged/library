package redis

import (
	"context"
	"time"

	"github.com/ArtEmerged/library/client/cache"
	"github.com/gomodule/redigo/redis"
)

var _ cache.Cache = (*client)(nil)

type handler func(ctx context.Context, conn redis.Conn) error

type Config interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type client struct {
	pool   *redis.Pool
	config Config
}

func NewClient(config Config) *client {
	pool := &redis.Pool{
		MaxIdle:     config.MaxIdle(),
		IdleTimeout: config.IdleTimeout(),
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialContext(ctx, "tcp", config.Address())
		},
	}

	return &client{
		pool:   pool,
		config: config,
	}
}
