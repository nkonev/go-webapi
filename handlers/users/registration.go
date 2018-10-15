package users

import (
	"github.com/nkonev/go-webapi/services"
	"time"
	"github.com/nkonev/go-webapi/models/token"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"github.com/nkonev/go-webapi/utils"
	"strings"
	"net/http"
	"errors"
	"github.com/nkonev/go-webapi/models/user"
)

type registrationHandler struct {
	mailer services.Mailer
	subject string
	bodyTemplate string
	url, confirmHandlerPath string
	confirmationTokenTtl time.Duration
	confirmationTokenModel token.ConfirmationRegistrationTokenModel
	userModel user.UserModel
}

func NewRegistrationHandler(mailer services.Mailer, subject string, bodyTemplate string,
	url, confirmHandlerPath string, confirmationTokenTtl time.Duration, confirmationTokenModel token.ConfirmationRegistrationTokenModel,
	userModel user.UserModel) *registrationHandler {
	return &registrationHandler{
		mailer: mailer,
		subject: subject,
		bodyTemplate: bodyTemplate,
		url: url,
		confirmHandlerPath: confirmHandlerPath,
		confirmationTokenTtl: confirmationTokenTtl,
		confirmationTokenModel: confirmationTokenModel,
		userModel: userModel,
	}
}

func (h *registrationHandler) Register(context echo.Context) error {
	d := &RegisterDTO{}
	if err := context.Bind(d); err != nil {
		return err
	}

	uuidStr := uuid.NewV4().String()
	link := generateConfirmLink(h.url, h.confirmHandlerPath, uuidStr)

	passwordHash, passwordHashErr := utils.HashPassword(d.Password)
	if passwordHashErr != nil {
		return passwordHashErr
	}

	if e := h.confirmationTokenModel.SaveTokenToRedis(uuidStr, &token.TempUser{d.Email, passwordHash}, h.confirmationTokenTtl); e != nil {
		return e
	}

	body := strings.Replace(h.bodyTemplate, "__link__", link, 1)
	h.mailer.SendMail(d.Email, h.subject, body)
	return context.JSON(http.StatusOK, utils.H{"message": "You successful registered, check your email"})
}

func generateConfirmLink(url, handlerPath string, uuid string) string {
	return url + handlerPath + "?token=" + uuid
}

func (h *registrationHandler) ConfirmRegistration(context echo.Context) error {
	confirmRegistrationToken := context.Request().URL.Query().Get("token")

	if len(confirmRegistrationToken) == 0 {
		return errors.New("Zero length token param")
	}
	if tempUser, err := h.confirmationTokenModel.GetValueByTokenFromRedis(confirmRegistrationToken); err != nil {
		return err
	} else {
		e := h.userModel.CreateUserByEmail(tempUser.Email, tempUser.PasswordHash)
		if e != nil {
			return e
		}
	}

	return context.JSON(http.StatusOK, utils.H{"message": "You successful confirm your registration"})
}
