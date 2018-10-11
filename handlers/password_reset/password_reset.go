package password_reset

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/nkonev/go-webapi/models/token"
	"github.com/nkonev/go-webapi/models/user"
	"github.com/nkonev/go-webapi/services"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

type handler struct {
	userModel user.UserModel
	passwordResetTokenModel token.PasswordResetTokenModel
	passwordResetTokenDuration time.Duration
	mailer services.Mailer
	passwordResetSubject, passwordResetBodyTemplate string
	url string
	handlerPath string
}

func NewHandler(passwordResetSubject, passwordResetBodyTemplate, url, handlerPath string, userModel user.UserModel, passwordResetTokenModel token.PasswordResetTokenModel, passwordResetTokenDuration time.Duration, mailer services.Mailer) *handler{
	return &handler{passwordResetSubject: passwordResetSubject, passwordResetBodyTemplate: passwordResetBodyTemplate, url: url, handlerPath: handlerPath, userModel: userModel, passwordResetTokenModel: passwordResetTokenModel, passwordResetTokenDuration: passwordResetTokenDuration, mailer: mailer}
}

type RequestPasswordResetDTO struct {
	Email       string
}

func (h *handler) RequestPasswordReset(c echo.Context) error {
	var dto = &RequestPasswordResetDTO{}
	if err := c.Bind(dto); err != nil {
		return err
	}

	user, userFindErr := h.userModel.FindByEmail(dto.Email)
	if userFindErr != nil {
		return userFindErr
	}
	if user != nil {
		maybeEmail := user.Email
		if !maybeEmail.Valid {
			log.Errorf("User have empty email was found. This shouldn' t happen.")
		} else {
			email := maybeEmail.String
			// generate password reset token
			uuidStr := uuid.NewV4().String()

			// save token to redis
			if err := h.passwordResetTokenModel.SaveTokenToRedis(uuidStr, h.passwordResetTokenDuration); err != nil {
				return err
			}

			// send this toke=n via email
			link := generateConfirmLink(h.url, h.handlerPath, uuidStr)
			body := strings.Replace(h.passwordResetBodyTemplate, "__link__", link, 1)
			h.mailer.SendMail(email, h.passwordResetSubject, body)
		}
	}
	// we just ignore case what de didn't find any user
	return nil
}

func generateConfirmLink(url, handlerPath string, uuid string) string {
	return url + handlerPath + "?token=" + uuid
}

func (h *handler) ConfirmPasswordReset(c echo.Context) error {
	return nil
}