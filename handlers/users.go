package users

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/go-echo-api-test-sample/models/user"
)

type resultLists struct {
	Users []user.User `json:"users"`
}

type handler struct {
	UserModel user.UserModel
}

type H map[string]interface{}

func NewHandler(u user.UserModel) *handler {
	return &handler{u}
}

type RegisterDTO struct {
	usernamed string
	password string
}

func (h *handler) GetIndex(c echo.Context) error {
	lists, e := h.UserModel.FindAll()
	if e != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": "We had error ("})
	}
	u := &resultLists{
		Users: lists,
	}
	return c.JSON(http.StatusOK, u)
}

func (h *handler) GetDetail(c echo.Context) error {
	id := c.Param("id")
	u, e := h.UserModel.FindByID(id)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": "We had error ("})
	}
	return c.JSON(http.StatusOK, u)
}

func (h *handler) GetProfile(c echo.Context) error {
	return c.JSON(http.StatusOK, H{"message": "You see your profile"})
}

func (h *handler) Register(context echo.Context) error {
	d := &RegisterDTO{}
	if err := context.Bind(d); err != nil {
		return err
	}
	return context.JSON(http.StatusOK, H{"message": "You successful registered"})
}