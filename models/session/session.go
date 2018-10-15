package session

import (
	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"time"
)

type SessionModel struct {
	Redis redis.Client
}

func getSessionKey(id string) string {
	return "session:" + id
}

func (sessionModel *SessionModel) CheckSession(key string) error {
	kv, e := sessionModel.Redis.HGetAll(getSessionKey(key)).Result()
	if e != nil {
		log.Errorf("Error during get session")
		return e
	}
	if len(kv) == 0 {
		return errors.Errorf("Got empty session from redis")
	}
	log.Infof("Successful checked session %v", kv)
	return nil
}

func (sessionModel *SessionModel) CreateSession(username string, sessionTtl time.Duration) (string, error) {
	sessionId := uuid.NewV4().String()
	log.Infof("Saving session %v with duration %v", sessionId, sessionTtl)

	if cmd := sessionModel.Redis.HSet(getSessionKey(sessionId), "login", username); cmd.Err() != nil {
		log.Errorf("Error during save session")
		return "error session", cmd.Err()
	}
	if err := sessionModel.Redis.Expire(getSessionKey(sessionId), sessionTtl).Err(); err != nil {
		log.Errorf("Error during set session expiration")
		return "error session", err
	}
	return sessionId, nil
}