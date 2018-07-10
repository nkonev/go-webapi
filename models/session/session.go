package session

import "github.com/go-redis/redis"

type SessionModel struct {
	Redis redis.Client
}