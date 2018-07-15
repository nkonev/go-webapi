package main

import "github.com/labstack/echo"
import (
	test "net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
	"io"
	"strings"
)


func request(method, path string, body io.Reader, e *echo.Echo) (int, string, http.Header) {
	req := test.NewRequest(method, path, body)
	Header := map[string][]string{
		echo.HeaderContentType: {"application/json"},
	}
	req.Header = Header
	rec := test.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String(), rec.HeaderMap
}

func TestUsers(t *testing.T) {
	e := configureEcho();
	defer e.Close()

	c, b, _ := request("GET", "/users", nil, e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestUser(t *testing.T) {
	e := configureEcho();
	defer e.Close()

	c, b, _ := request("GET", "/users/1", nil, e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestLoginSuccess(t *testing.T) {
	e := configureEcho();
	defer e.Close()

	_, _, hm := request("POST", "/auth2/login", strings.NewReader(`{"username": "root", "password": "password"}`), e)
	//assert.Equal(t, http.StatusOK, c)// todo
	assert.NotEmpty(t, hm.Get("Set-Cookie"))
	assert.Contains(t, hm.Get("Set-Cookie"), "SESSION=")
}

func TestLoginFail(t *testing.T) {
	e := configureEcho();
	defer e.Close()

	_, _, hm := request("POST", "/auth2/login", strings.NewReader(`{"username": "root", "password": "pass_-word"}`), e)
	//assert.Equal(t, http.StatusOK, c)// todo
	assert.Empty(t, hm.Get("Set-Cookie"))
}


func TestGetUsersWithSession(t *testing.T) {
	e := configureEcho();
	defer e.Close()
	// todo implement
}