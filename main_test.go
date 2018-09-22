package main

import "github.com/labstack/echo"
import (
	test "net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
	"io"
	"strings"
	"github.com/go-echo-api-test-sample/auth"
	serviceMocks "github.com/go-echo-api-test-sample/services/mocks"
	facebookMocks "github.com/go-echo-api-test-sample/handlers/facebook/mocks"
	"go.uber.org/dig"
	"github.com/go-echo-api-test-sample/services"
	"github.com/go-echo-api-test-sample/handlers/facebook"
	"github.com/go-echo-api-test-sample/db"
	"mvdan.cc/xurls"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	fb "github.com/huandu/facebook"
)


func request(method, path string, body io.Reader, e *echo.Echo, sessionCookie string) (int, string, http.Header) {
	req := test.NewRequest(method, path, body)
	Header := map[string][]string{
		echo.HeaderContentType: {"application/json"},
		echo.HeaderCookie: []string{constructSessionCookieHeaderValue(sessionCookie)},
	}
	req.Header = Header
	rec := test.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String(), rec.HeaderMap
}

func constructSessionCookieHeaderValue(session string) string {
	return auth.SESSION_COOKIE+"="+session
}

func getSession(headers http.Header) string {
	return strings.Replace(headers.Get(echo.HeaderSetCookie), auth.SESSION_COOKIE+"=", "", 1)
}

func mockMailer() services.Mailer{
	return &serviceMocks.Mailer{}
}

func mockFacebookClient() facebook.FacebookClient {
	return &facebookMocks.FacebookClient{}
}

func runTest(container *dig.Container, test func (e *echo.Echo)){
	if err := container.Invoke(func (e *echo.Echo){
		defer e.Close()

		test(e)
	}); err != nil {
		panic(err)
	}
}

func setUpContainerForIntegrationTests() *dig.Container {
	initViper()
	container := dig.New()
	container.Provide(db.ConfigureRedis)
	container.Provide(configureEcho)
	container.Provide(sessionModel)
	return container
}

func TestUsers(t *testing.T) {
	container := setUpContainerForIntegrationTests()
	container.Provide(mockMailer)
	container.Provide(mockFacebookClient)

	runTest(container, func (e *echo.Echo){
		c, b, _ := request("GET", "/users", nil, e, "")
		assert.Equal(t, http.StatusOK, c)
		assert.NotEmpty(t, b)
	})
}

func TestUser(t *testing.T) {
	container := setUpContainerForIntegrationTests()
	container.Provide(mockMailer)
	container.Provide(mockFacebookClient)

	runTest(container, func (e *echo.Echo){
		c, b, _ := request("GET", "/users/1", nil, e, "")
		assert.Equal(t, http.StatusOK, c)
		assert.NotEmpty(t, b)
	})
}

func TestLoginSuccess(t *testing.T) {
	container := setUpContainerForIntegrationTests()
	container.Provide(mockMailer)
	container.Provide(mockFacebookClient)

	runTest(container, func (e *echo.Echo){
		c, _, hm := request("POST", "/auth/login", strings.NewReader(`{"username": "root", "password": "password"}`), e, "")
		assert.Equal(t, http.StatusOK, c)
		assert.NotEmpty(t, hm.Get(echo.HeaderSetCookie))
		assert.Contains(t, hm.Get(echo.HeaderSetCookie), auth.SESSION_COOKIE+"=")

		session := getSession(hm)

		codeProfile, bodyProfile, _ := request("GET", "/profile", nil, e, session)
		assert.Equal(t, http.StatusOK, codeProfile)
		assert.Contains(t, bodyProfile, "You see your profile")
	})
}

func TestLoginFail(t *testing.T) {
	container := setUpContainerForIntegrationTests()
	container.Provide(mockMailer)
	container.Provide(mockFacebookClient)

	runTest(container, func (e *echo.Echo){
		c, _, hm := request("POST", "/auth/login", strings.NewReader(`{"username": "root", "password": "pass_-word"}`), e, "")
		assert.Equal(t, http.StatusInternalServerError, c)// todo
		assert.Empty(t, hm.Get("Set-Cookie"))
	})
}

func TestRegister(t *testing.T) {
	container := setUpContainerForIntegrationTests()
	container.Provide(mockFacebookClient)

	m := &serviceMocks.Mailer{}
	m.On("SendMail", "from@yandex.ru", "newroot@yandex.ru", "registration confirmation", mock.AnythingOfType("string"), "smtp.yandex.ru:465", "username", "password")
	container.Provide(func() services.Mailer {
		return m
	})
	container.Provide(mockFacebookClient)

	runTest(container, func (e *echo.Echo){
		c1, _, hm1 := request("POST", "/auth/register", strings.NewReader(`{"username": "newroot@yandex.ru", "password": "password"}`), e, "")
		assert.Equal(t, http.StatusOK, c1)
		assert.Empty(t, hm1.Get("Set-Cookie"))

		var emailBody string;
		emailBody = m.Calls[0].Arguments[3].(string)
		assert.Contains(t, emailBody, "http://localhost:1234/confirm/registration?token=")

		confirmUrl := xurls.Strict.FindString(emailBody)
		assert.Contains(t, confirmUrl, "http://localhost:1234/confirm/registration?token=")

		// confirm
		c2, _, _ := request("GET", confirmUrl, nil, e, "")
		assert.Equal(t, http.StatusOK, c2)

		// login
		c, _, hm := request("POST", "/auth/login", strings.NewReader(`{"username": "newroot@yandex.ru", "password": "password"}`), e, "")
		assert.Equal(t, http.StatusOK, c)
		assert.NotEmpty(t, hm.Get(echo.HeaderSetCookie))
		assert.Contains(t, hm.Get(echo.HeaderSetCookie), auth.SESSION_COOKIE+"=")


		m.AssertExpectations(t)
	})

}



func TestStaticIndex(t *testing.T) {

	container := setUpContainerForIntegrationTests()
	container.Provide(mockMailer)
	container.Provide(mockFacebookClient)

	runTest(container, func (e *echo.Echo){
		c, _, _ := request("GET", "/index.html", nil, e, "")
		assert.Equal(t, http.StatusMovedPermanently, c)
	})
}

func TestStaticRoot(t *testing.T) {

	container := setUpContainerForIntegrationTests()
	container.Provide(mockMailer)
	container.Provide(mockFacebookClient)

	runTest(container, func (e *echo.Echo){
		c, b, _ := request("GET", "/", nil, e, "")
		assert.Equal(t, http.StatusOK, c)
		assert.Contains(t, b, "Hello, world!")
	})
}


func TestStaticAssets(t *testing.T) {

	container := setUpContainerForIntegrationTests()
	container.Provide(mockMailer)
	container.Provide(mockFacebookClient)

	runTest(container, func (e *echo.Echo){
		c, b, _ := request("GET", "/assets/main.js", nil, e, "")
		assert.Equal(t, http.StatusOK, c)
		assert.Equal(t, `console.log("Hello world");`, b)
	})
}


func TestFacebookCallback(t *testing.T) {
	container := setUpContainerForIntegrationTests()
	container.Provide(mockMailer)

	f := &facebookMocks.FacebookClient{}
	f.On("Exchange",
		//mock.AnythingOfType("&oauth2.Config"),
		mock.Anything,
		mock.AnythingOfType("*context.emptyCtx"),
		"test0123").Return(&oauth2.Token{
		AccessToken: "accessToken456",
	}, nil)
	f.On("GetInfo", "accessToken456").Return(fb.Result{"email": "email@example.com"}, nil)

	container.Provide(mockMailer)
	container.Provide(func() facebook.FacebookClient {
		return f
	})

	runTest(container, func (e *echo.Echo){
		req := test.NewRequest("GET", "/auth/fb/callback?code=test0123", nil)
		header := map[string][]string{
			echo.HeaderContentType: {"application/json"},
		}
		req.Header = header
		rec := test.NewRecorder()
		e.ServeHTTP(rec, req)

		passedCode := f.Calls[0].Arguments[2]
		assert.Equal(t, "test0123", passedCode)

		f.AssertExpectations(t)
	})

}
