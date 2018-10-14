package password_reset

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/nkonev/go-webapi/models/token"
	"github.com/nkonev/go-webapi/models/user"
	"github.com/nkonev/go-webapi/services"
	"github.com/nkonev/go-webapi/utils"
	"github.com/satori/go.uuid"
	"net/http"
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
			if err := h.passwordResetTokenModel.SaveTokenToRedis(uuidStr, h.passwordResetTokenDuration, user.ID); err != nil {
				return err
			}

			// send this toke=n via email
			link := generateConfirmLink(h.url, h.handlerPath, uuidStr)
			body := strings.Replace(h.passwordResetBodyTemplate, "__link__", link, 1)
			h.mailer.SendMail(email, h.passwordResetSubject, body)
		}
	} else {
		// we just ignore case what de didn't find any user
		log.Infof("User with email '%v' not found. No email will send.", dto.Email)
	}
	return nil
}

func generateConfirmLink(url, handlerPath string, uuid string) string {
	return url + handlerPath + "?token=" + uuid
}

type ConfirmPasswordResetDto struct {
	PasswordResetToken string
	NewPassword string
}

func (h *handler) ConfirmPasswordReset(c echo.Context) error {
	d := &ConfirmPasswordResetDto{}
	if err := c.Bind(d); err != nil{
		return err
	}

	if passwordResetToken, err := h.passwordResetTokenModel.FindTokenInRedis(d.PasswordResetToken); err != nil {
		log.Infof("error during find password reset token in redis: '%v'", err)
		return c.JSON(http.StatusExpectationFailed, utils.H{"message": "Your password reset token is not found"})
	} else {
		passwordHash, passwordHashErr := utils.HashPassword(d.NewPassword)
		if passwordHashErr != nil {
			return passwordHashErr
		}

		if err := h.userModel.SetPassword(passwordResetToken, passwordHash); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, utils.H{"message": "You successfully changed your password"})
	}

	return nil
}