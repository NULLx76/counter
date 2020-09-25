package store

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisStore struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisStoreWithOptions(ctx context.Context, opts *redis.Options) *RedisStore {
	return &RedisStore{
		rdb: redis.NewClient(opts),
		ctx: ctx,
	}
}

func NewRedisStore(addr string) *RedisStore {
	return NewRedisStoreWithOptions(context.Background(),
		&redis.Options{
			Addr:        addr,
			MaxRetries:  5,
			DialTimeout: 5 * time.Second,
		},
	)
}

func (rs *RedisStore) set(key string, value Value) error {
	b, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	return rs.rdb.Set(rs.ctx, key, string(b), 0).Err()
}

func (rs *RedisStore) Create(key string, value Value) error {
	return rs.set(key, value)
}

func (rs *RedisStore) Delete(key string) error {
	return rs.rdb.Del(rs.ctx, key).Err()
}

func (rs *RedisStore) Get(key string) (Value, error) {
	n, err := rs.rdb.Exists(rs.ctx, key).Result()
	if err != nil {
		return Value{}, err
	} else if n == 0 {
		return Value{}, nil
	}

	val, err := rs.rdb.Get(rs.ctx, key).Result()
	if err != nil {
		return Value{}, err
	}

	var v Value
	if err := json.Unmarshal([]byte(val), &v); err != nil {
		return Value{}, err
	}

	return v, nil
}

func (rs *RedisStore) Increment(key string) error {
	val, err := rs.Get(key)
	if err != nil {
		return err
	}

	val.Count++

	return rs.set(key, val)
}

func (rs *RedisStore) Decrement(key string) error {
	val, err := rs.Get(key)
	if err != nil {
		return err
	}

	val.Count--

	return rs.set(key, val)
}

func (rs *RedisStore) Close() error {
	return rs.rdb.Close()
}
