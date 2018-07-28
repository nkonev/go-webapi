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
	"github.com/go-echo-api-test-sample/services/mocks"
	"github.com/stretchr/testify/mock"
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

func TestUsers(t *testing.T) {
	m := &mocks.Mailer{}
	e := configureEcho(m);
	defer e.Close()

	c, b, _ := request("GET", "/users", nil, e, "")
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestUser(t *testing.T) {
	m := &mocks.Mailer{}
	e := configureEcho(m);
	defer e.Close()

	c, b, _ := request("GET", "/users/1", nil, e, "")
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestLoginSuccess(t *testing.T) {
	m := &mocks.Mailer{}
	e := configureEcho(m);
	defer e.Close()

	c, _, hm := request("POST", "/auth2/login", strings.NewReader(`{"username": "root", "password": "password"}`), e, "")
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, hm.Get(echo.HeaderSetCookie))
	assert.Contains(t, hm.Get(echo.HeaderSetCookie), auth.SESSION_COOKIE+"=")

	session := getSession(hm)

	codeProfile, bodyProfile, _ := request("GET", "/profile", nil, e, session)
	assert.Equal(t, http.StatusOK, codeProfile)
	assert.Contains(t, bodyProfile, "You see your profile")
}

func TestLoginFail(t *testing.T) {
	m := &mocks.Mailer{}
	e := configureEcho(m);
	defer e.Close()

	c, _, hm := request("POST", "/auth2/login", strings.NewReader(`{"username": "root", "password": "pass_-word"}`), e, "")
	assert.Equal(t, http.StatusInternalServerError, c)// todo
	assert.Empty(t, hm.Get("Set-Cookie"))
}

func TestRegister(t *testing.T) {
	m := &mocks.Mailer{}
	m.On("SendMail", "from@yandex.ru", "root@yandex.ru", "registration confirmation", mock.AnythingOfType("string"), "smtp.yandex.ru:465", "username", "password")
	e := configureEcho(m);
	defer e.Close()

	c, _, hm := request("POST", "/auth2/register", strings.NewReader(`{"username": "root@yandex.ru", "password": "password"}`), e, "")
	assert.Equal(t, http.StatusOK, c)
	assert.Empty(t, hm.Get("Set-Cookie"))

	m.AssertExpectations(t)
}



func TestStaticIndex(t *testing.T) {
	m := &mocks.Mailer{}
	e := configureEcho(m);
	defer e.Close()

	c, _, _ := request("GET", "/index.html", nil, e, "")
	assert.Equal(t, http.StatusMovedPermanently, c)
}

func TestStaticRoot(t *testing.T) {
	m := &mocks.Mailer{}
	e := configureEcho(m);
	defer e.Close()

	c, b, _ := request("GET", "/", nil, e, "")
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "Hello, world!", b)
}

func TestStaticAssets(t *testing.T) {
	m := &mocks.Mailer{}
	e := configureEcho(m);
	defer e.Close()

	c, b, _ := request("GET", "/assets/main.js", nil, e, "")
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `console.log("Hello world");`, b)
}
