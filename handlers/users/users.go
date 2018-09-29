package users

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/nkonev/go-echo-api-test-sample/models/user"
	"github.com/nkonev/go-echo-api-test-sample/services"
	"github.com/satori/go.uuid"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"time"
	"github.com/jmoiron/sqlx"
	"errors"
	"github.com/nkonev/go-echo-api-test-sample/models/confirmation_token"
)

type resultLists struct {
	Users []user.User `json:"users"`
}

type handler struct {
	UserModel user.UserModel
}

type H map[string]interface{}

func NewHandler(u user.UserModel) *handler {
	return &handler{u}
}

type RegisterDTO struct {
	Username string // email
	Password string
}

func (h *handler) GetIndex(c echo.Context) error {
	lists, e := h.UserModel.FindAll()
	if e != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": "We had error ("})
	}
	u := &resultLists{
		Users: lists,
	}
	return c.JSON(http.StatusOK, u)
}

func (h *handler) GetDetail(c echo.Context) error {
	id := c.Param("id")
	u, e := h.UserModel.FindByID(id)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": "We had error ("})
	}
	return c.JSON(http.StatusOK, u)
}

func (h *handler) GetProfile(c echo.Context) error {
	return c.JSON(http.StatusOK, H{"message": "You see your profile"})
}

func (h *handler) Register(m services.Mailer, fromAdress string, subject string, bodyTemplate string,
	smtpHostPort string, smtpUserName string, smtpPassword string,
	url string, confirmationTokenTtl time.Duration, tm confirmation_token.ConfirmationTokenModel) echo.HandlerFunc {

	return func (context echo.Context) error {
		d := &RegisterDTO{}
		if err := context.Bind(d); err != nil {
			return err
		}

		uuidStr := uuid.NewV4().String()
		link := generateConfirmLink(url, uuidStr)

		passwordHash, passwordHashErr := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)
		if passwordHashErr != nil {
			return passwordHashErr
		}

		if e := tm.SaveTokenToRedis(uuidStr, &confirmation_token.TempUser{d.Username, string(passwordHash)}, confirmationTokenTtl); e != nil {
			return e
		}

		body := strings.Replace(bodyTemplate, "__link__", link, 1)
		m.SendMail(fromAdress, d.Username, subject, body, smtpHostPort, smtpUserName, smtpPassword)
		return context.JSON(http.StatusOK, H{"message": "You successful registered, check your email"})
	}
}

func generateConfirmLink(url string, uuid string) string {
	return url + "/confirm/registration?token="+uuid
}

func (h *handler) ConfirmRegistration(db *sqlx.DB, tm confirmation_token.ConfirmationTokenModel) echo.HandlerFunc {
	return func (context echo.Context) error {
		token := context.Request().URL.Query().Get("token")

		if len(token) == 0 {
			return errors.New("Zero length token param")
		}
		if tempUser, err := tm.GetValueByTokenFromRedis(token); err != nil {
			return err
		} else {
			e := h.UserModel.CreateUser(tempUser.Username, tempUser.PasswordHash)
			if e != nil {
				return e
			}
		}

		return context.JSON(http.StatusOK, H{"message": "You successful confirm your registration"})
	}
}
