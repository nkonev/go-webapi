package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/go-echo-api-test-sample/models/user"
)

type (
	UsersModelStub struct{}
)

func (u *UsersModelStub) FindByID(id string) (*user.User, error) {
	return &user.User{
		ID:   1,
		Name: "foo",
	}, nil
}
func (u *UsersModelStub) FindByLogin(id string) (*user.User, error) {
	return &user.User{
		ID:   1,
		Name: "foo",
	}, nil
}
func (u *UsersModelStub) FindAll() ([]user.User, error) {
	users := []user.User{}
	users = append(users, user.User{
		ID:   100,
		Name: "foo",
	})
	return users, nil
}
func (u *UsersModelStub)CreateUser(login string, passwordHash string) (error){
	return nil
}

func TestGetDetail(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	u := &UsersModelStub{}
	h := NewHandler(u)

	var userJSON = `{"id":1,"name":"foo","Surname":"","Lastname":""}`

	if assert.NoError(t, h.GetDetail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSON, rec.Body.String())
	}
}

func TestGetIndex(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users")

	u := &UsersModelStub{}
	h := NewHandler(u)

	var userJSON = `{"users":[{"id":100,"name":"foo","Surname":"","Lastname":""}]}`

	if assert.NoError(t, h.GetIndex(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSON, rec.Body.String())
	}
}
