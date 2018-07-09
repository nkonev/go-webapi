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
)

type MyServerStorer struct {
	Model user.UserModel
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
	log.Infof("Saving session Save")
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
			log.Infof("No cookie, creating session")
			return &MyClientState{Id: uuid.NewV4()}, nil
		}
		return nil, e
	}

	session0 := s.Model.Redis.HGetAll(c.Value)
	log.Infof("Loaded session %v", session0)
	return &MyClientState{set: *session0}, nil
}

// save session to redis
func (s *MySessionStorer) WriteState(w http.ResponseWriter, cstate authboss.ClientState, cse []authboss.ClientStateEvent) error {
	log.Infof("Saving session WriteState: %v | %v | %v", w, cstate, cse)

	// todo erase
	/*csrw, ok :=  w.(authboss.ClientStateResponseWriter)
	if ok {
		log.Infof("Got ClientStateResponseWriter %v", csrw)
	} else {
		log.Infof("Not ok")
	}*/

	//sss, ok := cstate.(MyClientState)

	for _, e := range cse {
		switch e.Kind {
		case authboss.ClientStateEventPut:
			s.Model.Redis.HSet("cstate.GetId().String()", e.Key, e.Value) // todo remove field
		case authboss.ClientStateEventDel:
			s.Model.Redis.HDel(e.Key)
		}
	}

	return nil
}


type MyClientState struct {
	Id uuid.UUID
	set redis.StringStringMapCmd
}

func (s *MyClientState) Get(key string) (string, bool) {
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

func (s *MyClientState) GetId() uuid.UUID {
	return s.Id
}