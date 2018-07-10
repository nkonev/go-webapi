package auth

import (
	"github.com/volatiletech/authboss"
	"context"
	"github.com/labstack/gommon/log"
	"github.com/go-echo-api-test-sample/models"
	"net/http"
	"github.com/go-echo-api-test-sample/models/session"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"github.com/pkg/errors"
)

type MyServerStorer struct {
	Model user.UserModelImpl
}

var session_cookie = "SESSION"

func (s *MyServerStorer) Load(ctx context.Context, key string) (authboss.User, error)  {
	log.Infof("Loading user")
	uu, e  := s.Model.FindByLogin(key)
	if e != nil {
		return nil, e
	}
	if uu == nil {
		return nil, authboss.ErrUserNotFound
	}
	return uu, nil
}

func (s *MyServerStorer) Save(ctx context.Context, user authboss.User) (error)  {
	log.Infof("Saving user")
	return nil // todo implement
}


type MySessionStorer struct {
	Model session.SessionModel
}

// read session from redis with cookie
func (s *MySessionStorer) ReadState(r *http.Request) (authboss.ClientState, error) {
	c, e := r.Cookie(session_cookie)
	if e != nil {
		if e == http.ErrNoCookie {
			u := uuid.NewV4().String()
			log.Infof("No cookie, creating session %v", u)
			ss := MyClientStateImpl{Id: u}
			return &ss, nil
		}
		return nil, e
	}

	session0 := s.Model.Redis.HGetAll(c.Value)
	log.Infof("Loaded session %v", session0)
	return &MyClientStateImpl{Id: c.Value, set: *session0}, nil
}

// save session to redis
func (s *MySessionStorer) WriteState(w http.ResponseWriter, cstate authboss.ClientState, cse []authboss.ClientStateEvent) error {


	m, ok := cstate.(MyClientState)
	if !ok {
		return errors.Errorf("Cannot cast to MyClientState")
	}
	sessionId := m.GetSessionId()
	log.Infof("Saving session %v", sessionId)

	for _, e := range cse {
		switch e.Kind {
		case authboss.ClientStateEventPut:
			s.Model.Redis.HSet(sessionId, e.Key, e.Value)
		case authboss.ClientStateEventDel:
			s.Model.Redis.HDel(sessionId, e.Key)
		}
	}

	return nil
}

type MyClientState interface {
	authboss.ClientState
	GetSessionId() string
}

type MyClientStateImpl struct {
	Id string
	set redis.StringStringMapCmd
}

func (s MyClientStateImpl) Get(key string) (string, bool) {
	map0, err := s.set.Result()
	if err != nil {
		log.Errorf("Has error during get value from HSET %v", err)
		return "", false
	}
	v := map0[key]
	if v == "" {
		return "", false
	} else {
		return v, false
	}
}

func (s MyClientStateImpl) GetSessionId() string {
	return s.Id
}