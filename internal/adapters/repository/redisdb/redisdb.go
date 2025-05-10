package redisdb

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Collections names
const (
	ProductCollection = "products"
	OrderCollection   = "orders"
	UserCollection    = "users"
	AddressCollection = "addresses"
	SettingCollection = "settings"
)

// NewRedisConnection establishes a connection to Redis
func NewRedisConnection(uri string) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})

	// Ping the database to verify connection
	status := client.Ping(ctx)
	if status.Err() != nil {
		return nil, status.Err()
	}
	// Set a timeout for Redis operations

	return client, nil
}
