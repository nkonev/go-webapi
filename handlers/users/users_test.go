package users

import (
	"github.com/guregu/null"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/nkonev/go-webapi/models/user"
	"github.com/stretchr/testify/assert"
)

type (
	UsersModelStub struct{}
)

func (u *UsersModelStub) FindByID(id int) (*user.User, error) {
	return &user.User{
		ID:    1,
		Email: null.StringFrom("foo"),
	}, nil
}
func (u *UsersModelStub) FindByEmail(email string) (*user.User, error) {
	return &user.User{
		ID:    1,
		Email: null.StringFrom("foo"),
	}, nil
}
func (u *UsersModelStub) FindAll() ([]user.User, error) {
	users := []user.User{}
	users = append(users, user.User{
		ID:    100,
		Email: null.StringFrom("foo"),
	})
	return users, nil
}
func (u *UsersModelStub) CreateUserByEmail(login string, passwordHash string) error {
	return nil
}
func (u *UsersModelStub) CreateUserByFacebook(facebookId string) error {
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

	var userJSON = `{"id":1,"email":"foo","creationType":"","facebookId":null}`

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

	var userJSON = `{"users":[{"id":100,"email":"foo","creationType":"","facebookId":null}]}`

	if assert.NoError(t, h.GetIndex(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSON, rec.Body.String())
	}
}
