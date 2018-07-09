package session

import "github.com/go-redis/redis"

type Session struct {
	id string
}

type SessionModel struct {
	Redis redis.Client
}