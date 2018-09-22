package db

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
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

func ConfigureRedis()  *redis.Client {
	redisAddr := viper.GetString("redis.addr")
	redisPassword := viper.GetString("redis.password")
	redisDbNum := viper.GetInt("redis.db")
	redisFlushOnStart := viper.GetBool("redis.flushOnStart")

	return ConnectRedis(redisAddr, redisPassword, redisDbNum, redisFlushOnStart)
}