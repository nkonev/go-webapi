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


func request(method, path string, body io.Reader, e *echo.Echo) (int, string) {
	req := test.NewRequest(method, path, body)
	rec := test.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestUsers(t *testing.T) {
	e := configureEcho();
	defer e.Close()

	c, b := request("GET", "/users", nil, e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestUser(t *testing.T) {
	e := configureEcho();
	defer e.Close()

	c, b := request("GET", "/users/1", nil, e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestLogin(t *testing.T) {
	e := configureEcho();
	defer e.Close()

	c, b := request("POST", "/auth2/login", strings.NewReader(`{"username": "nick", "password": "lol"}`), e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestGetUsersWithSession(t *testing.T) {
	e := configureEcho();
	defer e.Close()
	// todo implement
}