package auth

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/go-echo-api-test-sample/models/session"
	"net/http"
	"regexp"
	"golang.org/x/crypto/bcrypt"
	"github.com/pkg/errors"
	"github.com/go-echo-api-test-sample/models/user"
	"github.com/satori/go.uuid"
	"time"
)

const SESSION_COOKIE  = "SESSION";

func CheckSession(context echo.Context, next echo.HandlerFunc, sessionModel session.SessionModel, whitelist []string) error {
	c, e := context.Request().Cookie(SESSION_COOKIE)
	if e != nil {
		if e == http.ErrNoCookie{
			for _, regexp0 := range whitelist {
				r, _ := regexp.Compile(regexp0) //todo optimize
				if r.MatchString(context.Request().RequestURI) {
					log.Infof("Skipping authentication for %v", regexp0)
					return next(context)
				}
			}
		}
		log.Errorf("Error get %v cookie", SESSION_COOKIE)
		return e
	}

	kv, e := sessionModel.Redis.HGetAll(c.Value).Result()
	if e != nil {
		log.Errorf("Error during get session")
		return e
	}
	if len(kv) == 0 {
		return errors.Errorf("Got empty session from redis")
	}
	log.Infof("Loaded session %v", kv)

	return next(context)
}

type LoginDTO struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginManager(context echo.Context, sessionModel session.SessionModel, userModel user.UserModel) error {
	m := new(LoginDTO)
	if err := context.Bind(m); err != nil {
		return err
	}

	userEntity, e  := userModel.FindByLogin(m.Username)
	if e != nil {
		return e
	}
	if userEntity == nil {
		return errors.Errorf("User %v not found", m.Username)
	}


	ep := bcrypt.CompareHashAndPassword([]byte(userEntity.GetPassword()), []byte(m.Password))
	if ep != nil {
		return errors.Errorf("Bad password")
	}

	sessionId := uuid.NewV4().String()
	ttl := "30m"
	log.Infof("Saving session %v with duration %v", sessionId, ttl)
	cmd := sessionModel.Redis.HSet(sessionId, "login", m.Username)
	d, _ := time.ParseDuration(ttl)
	sessionModel.Redis.Expire(sessionId, d)
	if cmd.Err() != nil {
		log.Errorf("Error during save session")
		return cmd.Err()
	}

	c := &http.Cookie{
		Name: SESSION_COOKIE,
		Value: sessionId,
	}
	http.SetCookie(context.Response(), c)

	return nil
}