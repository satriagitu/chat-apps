package repository_test

import (
	"context"
	"testing"

	"chat-apps/internal/repository"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_Set(t *testing.T) {
	ctx := context.Background()

	// Create a mock Redis client and mock controller
	mock, mockCtrl := redismock.NewClientMock()

	// Set up the expected behavior on the mock client
	key := "test_key"
	value := "test_value"
	mockCtrl.ExpectSet(key, value, 0).SetVal("OK")

	// Pass the mock client to the repository
	r := repository.NewRedisRepository(mock, ctx)

	// Call the Set method
	err := r.Set(ctx, key, value)

	// Validate the results
	assert.NoError(t, err)

	// Ensure all expectations are met
	err = mockCtrl.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRedisRepository_Get(t *testing.T) {
	ctx := context.Background()

	// Create a mock Redis client and mock controller
	mock, mockCtrl := redismock.NewClientMock()

	// Set up the expected behavior on the mock client
	key := "test_key"
	expectedValue := "test_value"
	mockCtrl.ExpectGet(key).SetVal(expectedValue)

	// Pass the mock client to the repository
	r := repository.NewRedisRepository(mock, ctx)

	// Call the Get method
	actualValue, err := r.Get(ctx, key)

	// Validate the results
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, actualValue)
	t.Log("actualValue:", actualValue)
	// Ensure all expectations are met
	err = mockCtrl.ExpectationsWereMet()
	assert.NoError(t, err)
}
