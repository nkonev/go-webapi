package auth

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/nkonev/go-webapi/models/session"
	"net/http"
	"regexp"
	"golang.org/x/crypto/bcrypt"
	"github.com/pkg/errors"
	"github.com/nkonev/go-webapi/models/user"
	"time"
)

const SESSION_COOKIE = "SESSION";

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

	if err := sessionModel.CheckSession(c.Value); err!= nil {
		return err
	}

	return next(context)
}

type LoginDTO struct{
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginManager(context echo.Context, sessionModel session.SessionModel, userModel user.UserModel, sessionTtl time.Duration) error {
	m := new(LoginDTO)
	if err := context.Bind(m); err != nil {
		return err
	}

	userEntity, e := userModel.FindByEmail(m.Email)
	if e != nil {
		return e
	}
	if userEntity == nil {
		return errors.Errorf("User %v not found", m.Email)
	}

	if !userEntity.Password.Valid {
		log.Errorf("Null password in database")
		return errors.Errorf("Bad password")
	}
	passwordCompareError := bcrypt.CompareHashAndPassword([]byte(userEntity.Password.String), []byte(m.Password))
	if passwordCompareError != nil {
		return errors.Errorf("Bad password")
	}

	sessionId, sessionCreateError := sessionModel.CreateSession(m.Email, sessionTtl)
	if sessionCreateError != nil {
		return sessionCreateError
	}

	c := &http.Cookie{
		Name: SESSION_COOKIE,
		Value: sessionId,
	}
	http.SetCookie(context.Response(), c)

	return nil
}