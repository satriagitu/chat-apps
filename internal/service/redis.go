package service

import (
	"chat-apps/internal/repository"
	"context"
)

type CacheService interface {
	CacheData(ctx context.Context, key string, value string) error
	RetrieveData(ctx context.Context, key string) (string, error)
}

type cacheService struct {
	cache repository.RedisRepository
}

func NewCacheService(cache repository.RedisRepository) CacheService {
	return &cacheService{cache: cache}
}

func (u *cacheService) CacheData(ctx context.Context, key string, value string) error {
	return u.cache.Set(ctx, key, value)
}

func (u *cacheService) RetrieveData(ctx context.Context, key string) (string, error) {
	return u.cache.Get(ctx, key)
}
