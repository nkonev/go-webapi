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

func checkUrlInWhitelist(whitelist []regexp.Regexp, uri string) bool {
	for _, regexp0 := range whitelist {
		if regexp0.MatchString(uri) {
			log.Infof("Skipping authentication for %v because it matches %v", uri, regexp0.String())
			return true
		}
	}
	return false
}

func CheckSession(context echo.Context, next echo.HandlerFunc, sessionModel session.SessionModel, whitelist []regexp.Regexp) error {
	if checkUrlInWhitelist(whitelist, context.Request().RequestURI){
		return next(context)
	}

	c, e := context.Request().Cookie(SESSION_COOKIE)
	if e != nil {
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

func LoginManager(context echo.Context, sessionModel session.SessionModel, userModel user.UserModel, sessionTtl time.Duration) error {
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
	log.Infof("Saving session %v with duration %v", sessionId, sessionTtl)

	if cmd := sessionModel.Redis.HSet(sessionId, "login", m.Username); cmd.Err() != nil {
		log.Errorf("Error during save session")
		return cmd.Err()
	}
	if err := sessionModel.Redis.Expire(sessionId, sessionTtl).Err(); err != nil {
		log.Errorf("Error during set session expiration")
		return err
	}

	c := &http.Cookie{
		Name: SESSION_COOKIE,
		Value: sessionId,
	}
	http.SetCookie(context.Response(), c)

	return nil
}