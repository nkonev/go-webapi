package users

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/nkonev/go-webapi/models/token"
	"github.com/nkonev/go-webapi/models/user"
	"github.com/nkonev/go-webapi/services"
	"github.com/nkonev/go-webapi/utils"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type resultLists struct {
	Users []user.User `json:"users"`
}

type handler struct {
	UserModel user.UserModel
}

func NewHandler(u user.UserModel) *handler {
	return &handler{u}
}

type RegisterDTO struct {
	Email    string // email
	Password string
}

func (h *handler) GetIndex(c echo.Context) error {
	lists, e := h.UserModel.FindAll()
	if e != nil {
		return c.JSON(http.StatusInternalServerError, utils.H{"error": "We had error ("})
	}
	u := &resultLists{
		Users: lists,
	}
	return c.JSON(http.StatusOK, u)
}

func (h *handler) GetDetail(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	u, e := h.UserModel.FindByID(idInt)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, utils.H{"error": "We had error ("})
	}
	return c.JSON(http.StatusOK, u)
}

func (h *handler) GetProfile(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.H{"message": "You see your profile"})
}

func (h *handler) Register(m services.Mailer, subject string, bodyTemplate string,
	url, confirmHandlerPath string, confirmationTokenTtl time.Duration, tm token.ConfirmationRegistrationTokenModel) echo.HandlerFunc {

	return func(context echo.Context) error {
		d := &RegisterDTO{}
		if err := context.Bind(d); err != nil {
			return err
		}

		uuidStr := uuid.NewV4().String()
		link := generateConfirmLink(url, confirmHandlerPath, uuidStr)

		passwordHash, passwordHashErr := utils.HashPassword(d.Password)
		if passwordHashErr != nil {
			return passwordHashErr
		}

		if e := tm.SaveTokenToRedis(uuidStr, &token.TempUser{d.Email, passwordHash}, confirmationTokenTtl); e != nil {
			return e
		}

		body := strings.Replace(bodyTemplate, "__link__", link, 1)
		m.SendMail(d.Email, subject, body)
		return context.JSON(http.StatusOK, utils.H{"message": "You successful registered, check your email"})
	}
}

func generateConfirmLink(url, handlerPath string, uuid string) string {
	return url + handlerPath + "?token=" + uuid
}

func (h *handler) ConfirmRegistration(tm token.ConfirmationRegistrationTokenModel) echo.HandlerFunc {
	return func(context echo.Context) error {
		confirmRegistrationToken := context.Request().URL.Query().Get("token")

		if len(confirmRegistrationToken) == 0 {
			return errors.New("Zero length token param")
		}
		if tempUser, err := tm.GetValueByTokenFromRedis(confirmRegistrationToken); err != nil {
			return err
		} else {
			e := h.UserModel.CreateUserByEmail(tempUser.Email, tempUser.PasswordHash)
			if e != nil {
				return e
			}
		}

		return context.JSON(http.StatusOK, utils.H{"message": "You successful confirm your registration"})
	}
}
