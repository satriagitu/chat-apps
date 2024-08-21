package service_test

import (
	"chat-apps/internal/repository/mocks"
	"chat-apps/internal/service"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheData(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	redisService := service.NewCacheService(mockRedisRepo)

	ctx := context.Background()
	mockRedisRepo.On("Set", ctx, "mykey", "myvalue").Return(nil)
	err := redisService.CacheData(ctx, "mykey", "myvalue")

	assert.Nil(t, err)
}

func TestRetrieveData(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	redisService := service.NewCacheService(mockRedisRepo)

	ctx := context.Background()
	expectedValue := "expectedValue"
	mockRedisRepo.On("Get", ctx, "mykey").Return(expectedValue, nil)
	value, err := redisService.RetrieveData(ctx, "mykey")

	assert.Nil(t, err)
	assert.Equal(t, expectedValue, value)
}
