package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) // ✅ add this
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, password string, db int) Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &redisClient{
		client: rdb,
	}
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisClient) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisClient) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return r.client.Scan(ctx, cursor, match, count).Result()
}
