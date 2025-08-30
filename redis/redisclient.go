package redisclient

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

// Init initializes the Redis client (only once)
func Init(addr, password string, db int) {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
			Protocol: 2,
		})
	})
}

// Get returns the Redis client instance
func Get() *redis.Client {
	if client == nil {
		panic("Redis client not initialized. Call redisclient.Init() first.")
	}
	return client
}
