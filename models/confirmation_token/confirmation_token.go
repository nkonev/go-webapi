package confirmation_token

import (
	"time"
	"github.com/go-redis/redis"
)

type ConfirmationTokenModel interface {
	SaveTokenToRedis(token string, u *TempUser, confirmationTokenTtl time.Duration) error
	GetValueByTokenFromRedis(token string) (TempUser, error)
}

type TempUser struct {
	Email        string
	PasswordHash string
}

type confirmationTokenModelImpl struct {
	redis redis.Client
}

func NewConfirmationTokenModel(redis *redis.Client) ConfirmationTokenModel {
	return &confirmationTokenModelImpl{redis: *redis}
}

const fieldUserName = "username"
const fieldPassword = "password"

func getKey(token string) string {
	return "registration:"+token;
}

func (i *confirmationTokenModelImpl) SaveTokenToRedis(token string, u *TempUser, confirmationTokenTtl time.Duration) error {
	userData := map[string]interface{}{
		fieldUserName: u.Email,
		fieldPassword: u.PasswordHash,
	}
	c := i.redis.HMSet(getKey(token), userData)
	if c.Err() != nil {
		return c.Err()
	}
	i.redis.Expire(getKey(token), confirmationTokenTtl)
	return nil

}

func (i *confirmationTokenModelImpl) GetValueByTokenFromRedis(token string) (TempUser, error) {
	redisResponse := i.redis.HGetAll(getKey(token))
	if map0, err := redisResponse.Result(); err != nil {
		return TempUser{}, redisResponse.Err()
	} else {
		username := map0[fieldUserName]
		password := map0[fieldPassword]

		return TempUser{Email: username, PasswordHash:password}, nil
	}

}