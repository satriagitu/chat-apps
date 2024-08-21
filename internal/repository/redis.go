package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisRepository interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
}

type redisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(client *redis.Client, ctx context.Context) RedisRepository {
	return &redisRepository{client: client, ctx: ctx}
}

func (r *redisRepository) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
