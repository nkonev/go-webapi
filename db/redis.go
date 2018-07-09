package db

import (
	"github.com/go-redis/redis"
)

func Connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "172.24.0.3:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}
