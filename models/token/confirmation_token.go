package token

import (
	"time"
	"github.com/go-redis/redis"
)

type ConfirmationRegistrationTokenModel interface {
	SaveTokenToRedis(token string, u *TempUser, confirmationTokenTtl time.Duration) error
	GetValueByTokenFromRedis(token string) (TempUser, error)
	FindTokenByEmail(email string) (string, bool, error)
	DeleteToken(token string) error
}

type TempUser struct {
	Email        string
	PasswordHash string
}

type confirmationTokenModelImpl struct {
	redis redis.Client
}

func NewConfirmationTokenModel(redis *redis.Client) ConfirmationRegistrationTokenModel {
	return &confirmationTokenModelImpl{redis: *redis}
}

const fieldUserName = "username"
const fieldPassword = "password"
const RegistrationTokenPrefix = "registration:token:"
func getKey(token string) string {
	return RegistrationTokenPrefix +token;
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

func (i *confirmationTokenModelImpl) DeleteToken(token string) error {
	return i.redis.Del(getKey(token)).Err()
}

func (i *confirmationTokenModelImpl) FindTokenByEmail(email string) (string, bool, error) {
	iter := i.redis.Scan(0, RegistrationTokenPrefix+"*", 128).Iterator()
	for iter.Next() {
		key := iter.Val()
		if found, err := i.findTokenMatchingEmail(key, email); err != nil {
			return "", false, err
		} else if found {
			return key[len(RegistrationTokenPrefix):], true, nil
		} // else continue
	}
	if err := iter.Err(); err != nil {
		return "", false, err
	} else {
		return "", false, nil // not found
	}
}

func (impl *confirmationTokenModelImpl) findTokenMatchingEmail(key string, email string) (bool, error) {
	iter := impl.redis.HScan(key, 0, fieldUserName, 8).Iterator()
	for iter.Next() {
		if iter.Val() == email {
			return true, nil
		} // else continue
	}
	if err := iter.Err(); err != nil {
		return false, err
	} else {
		return false, nil // not found
	}
}