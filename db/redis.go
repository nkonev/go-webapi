package db

import (
	"github.com/go-redis/redis"
)

func ConnectRedis(redisAddr string, redisPassword string, db int, flushOnStart bool) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // no password set
		DB:       db,  // use default DB
	})
	if flushOnStart {
		client.FlushDB()
	}
	return client
}
