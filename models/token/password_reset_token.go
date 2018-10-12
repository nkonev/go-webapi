package token

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type PasswordResetTokenModel interface {
	SaveTokenToRedis(token string, passwordResetTokenTtl time.Duration, userId int) error
	FindTokenInRedis(token string) (int, error)
}

type passwordResetTokenModelImpl struct {
	redis redis.Client
}

func NewPasswordResetTokenModel(redis *redis.Client) PasswordResetTokenModel {
	return &passwordResetTokenModelImpl{redis: *redis}
}

func getPasswordResetKey(token string) string {
	return "password:reset:token:" + token
}

func (model *passwordResetTokenModelImpl) SaveTokenToRedis(token string, passwordResetTokenTtl time.Duration, userId int) error {
	return model.redis.Set(getPasswordResetKey(token), userId, passwordResetTokenTtl).Err()
}

func (model *passwordResetTokenModelImpl) FindTokenInRedis(token string) (int, error) {
	getResult := model.redis.Get(getPasswordResetKey(token))
	if getResult.Err() != nil {
		return -1, getResult.Err()
	}
	return strconv.Atoi(getResult.Val())
}