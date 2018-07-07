package auth

import (
	"github.com/volatiletech/authboss"
	"context"
	"github.com/labstack/gommon/log"
	"github.com/go-echo-api-test-sample/models"
)

type MyServerStorer struct {
	Model user.UserModel
}

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
	return nil // todo implement
}

