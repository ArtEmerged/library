package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

func (c *client) Get(ctx context.Context, key string, in interface{}) error {
	var value interface{}
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		var errEx error
		value, errEx = conn.Do("GET", key)
		if errEx != nil {
			return errEx
		}

		return nil
	})
	if err != nil {
		return err
	}

	data, err := redis.Bytes(value, err)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal value")
	}

	err = json.Unmarshal(data, in)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal value")
	}

	return nil
}

func (c *client) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "can't marshal value")
	}

	err = c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, err := conn.Do("SET", redis.Args{key}.Add(data)...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) HSet(ctx context.Context, key, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "can't marshal value")
	}

	err = c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, err := conn.Do("HSET", redis.Args{key, field}.Add(data)...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) HGet(ctx context.Context, key, field string, in interface{}) error {
	var value interface{}
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		var errEx error
		value, errEx = conn.Do("HGET", key, field)
		if errEx != nil {
			return errEx
		}

		return nil
	})
	if err != nil {
		return err
	}

	data, err := redis.Bytes(value, err)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal value")
	}

	err = json.Unmarshal(data, in)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal value")
	}

	return nil
}

func (c *client) HDel(ctx context.Context, key string, field ...string) error {
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		args := redis.Args{key}
		for _, f := range field {
			args = args.Add(f)
		}

		_, err := conn.Do("HDEL", args...)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// TODO need to fix
func (c *client) HGetAll(ctx context.Context, key string) ([]interface{}, error) {
	var values []interface{}
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		var errEx error
		values, errEx = redis.Values(conn.Do("HGETALL", key))
		if errEx != nil {
			return errEx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (c *client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, err := conn.Do("EXPIRE", key, int(expiration.Seconds()))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Ping(ctx context.Context) error {
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, err := conn.Do("PING")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Close() error {
	return c.pool.Close()
}

func (c *client) execute(ctx context.Context, handler handler) error {
	conn, err := c.getConnect(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("failed to close redis connection: %v\n", err)
		}
	}()

	err = handler(ctx, conn)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) getConnect(ctx context.Context) (redis.Conn, error) {
	getConnTimeoutCtx, cancel := context.WithTimeout(ctx, c.config.ConnectionTimeout())
	defer cancel()

	conn, err := c.pool.GetContext(getConnTimeoutCtx)
	if err != nil {
		log.Printf("failed to get redis connection: %v\n", err)

		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}
