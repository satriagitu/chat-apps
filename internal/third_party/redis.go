package third_party

import (
	"github.com/go-redis/redis/v8"
)

// InitRedis initializes and returns a Redis client and context
func InitRedis() *redis.Client {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address// Tambahkan opsi lain jika diperlukan, seperti password, DB, dll.
	})

	return rdb
}
