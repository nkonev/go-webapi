package auth

import (
	"github.com/volatiletech/authboss"
	"context"
	"github.com/labstack/gommon/log"
	"github.com/go-echo-api-test-sample/models/user"
	"net/http"
	"github.com/go-echo-api-test-sample/models/session"
	"github.com/satori/go.uuid"
	"github.com/pkg/errors"
)

type MyUserStorer struct {
	Model user.UserModelImpl
}

const session_cookie = "SESSION"

func (s *MyUserStorer) Load(ctx context.Context, key string) (authboss.User, error)  {
	log.Infof("Try to find user '%v'", key)
	uu, e  := s.Model.FindByLogin(key)
	if e != nil {
		return nil, e
	}
	if uu == nil {
		return nil, authboss.ErrUserNotFound
	}
	return uu, nil
}

func (s *MyUserStorer) Save(ctx context.Context, user authboss.User) (error)  {
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
			log.Infof("No cookie named %v, creating session %v", session_cookie, u)
			ss := MyClientSessionImpl{id: u}
			return &ss, nil
		}
		return nil, e
	}

	kv, e := s.Model.Redis.HGetAll(c.Value).Result()
	if e != nil {
		log.Panicf("Cannot deserialize map %v", e)
	}
	log.Infof("Loaded session %v", kv)
	return &MyClientSessionImpl{id: c.Value, kv: kv}, nil
}

// save session to redis
func (s *MySessionStorer) WriteState(w http.ResponseWriter, cstate authboss.ClientState, cse []authboss.ClientStateEvent) error {

	m, ok := cstate.(MyClientSession)
	if !ok {
		return errors.Errorf("Cannot cast to MyClientSession")
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

	c := &http.Cookie{
		Name: session_cookie,
		Value: sessionId,
	}
	http.SetCookie(w, c)

	return nil
}

type MyClientSession interface {
	authboss.ClientState

	GetSessionId() string
}

type MyClientSessionImpl struct {
	id string
	kv map[string]string
}

func (s MyClientSessionImpl) Get(key string) (string, bool) {
	v := s.kv[key]
	if v == "" {
		return "", false
	} else {
		return v, true
	}
}

func (s MyClientSessionImpl) GetSessionId() string {
	return s.id
}