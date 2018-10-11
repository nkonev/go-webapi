package token

import (
	"github.com/go-redis/redis"
	"time"
)

type PasswordResetTokenModel interface {
	SaveTokenToRedis(token string, passwordResetTokenTtl time.Duration) error
	HasTokenInRedis(token string) (bool, error)
}

func NewPasswordResetTokenModel(redis redis.Client) *passwordResetTokenModelImpl {
	return &passwordResetTokenModelImpl{redis: redis}
}

type passwordResetTokenModelImpl struct {
	redis redis.Client
}

func (model *passwordResetTokenModelImpl) SaveTokenToRedis(token string, passwordResetTokenTtl time.Duration) error {
	return nil // todo implement
}

func (model *passwordResetTokenModelImpl) HasTokenInRedis(token string) (bool, error) {
	return false, nil //todo implement
}