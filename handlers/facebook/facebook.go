package facebook

import (
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"net/http"
	"golang.org/x/oauth2/facebook"
	"github.com/labstack/gommon/log"
	"errors"
	fb "github.com/huandu/facebook"
	"github.com/go-echo-api-test-sample/models/user"
)

type handler struct {
	config *oauth2.Config
	userModel *user.UserModelImpl
}

func NewHandler(facebookClientId, facebookSecret, urlFbCallback string, userModel *user.UserModelImpl) *handler {
	if len(facebookClientId) == 0 {
		log.Panicf("facebookClientId is empty")
	}
	if len(facebookSecret) == 0 {
		log.Panicf("facebookSecret is empty")
	}
	if len(urlFbCallback) == 0 {
		log.Panicf("urlFbCallback is empty")
	}

	return &handler{
		userModel: userModel,
		config: &oauth2.Config{
			ClientID:     facebookClientId,
			ClientSecret: facebookSecret,
			RedirectURL:  urlFbCallback, // Url that handles the response given
			// after authentication from the facebook.
			Scopes:   []string{"email"}, // The permission you want for the app
			Endpoint: facebook.Endpoint,
		}}
}

func (h *handler) RedirectForLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		url := h.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(c.Response().Writer, c.Request(), url, 301) // This will redirect to the facebook login page for
		return nil
	}
}

func (h *handler) CallBackHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.Request().URL.Query().Get("code")
		if len(code) == 0 {
			return errors.New("zero length code is not allowed")
		}
		log.Infof("You are successfully get code %v", code)
		// Handle it on your way
		token, err := h.config.Exchange(c.Request().Context(), code)
		if err != nil {
			return err
		}
		log.Infof("Successfully get access token %v which will expired at %v", token.AccessToken, token.Expiry)

		res, err := fb.Get("/me", fb.Params{
			"fields": "name,email,id",
			"access_token": token.AccessToken,
		})
		if err != nil {
			return err
		}
		log.Infof("Got facebook response: %v", res)

		emailRaw := res.Get("email")
		log.Infof("Got facebook email: %v", emailRaw)
		email := emailRaw.(string)

		if err := h.userModel.CreateUser(email, ""); err == nil {
			c.HTML(200, `You are successfully registered as `+email)
			return nil
		} else {
			return err
		}
	}
}