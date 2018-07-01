package main

import "github.com/labstack/echo"
import (
	test "net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
)


func request(method, path string, e *echo.Echo) (int, string) {
	req := test.NewRequest(method, path, nil)
	rec := test.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestUsers(t *testing.T) {
	e := configureEcho();
	defer e.Close()

	c, b := request("GET", "/users", e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestUser(t *testing.T) {
	e := configureEcho();
	defer e.Close()

	c, b := request("GET", "/users/1", e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}