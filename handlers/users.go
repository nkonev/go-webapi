package users

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/go-echo-api-test-sample/models"
)

type resultLists struct {
	Users []user.User `json:"users"`
}

type handler struct {
	UserModel user.UserModelImpl
}

type H map[string]interface{}

func NewHandler(u user.UserModelImpl) *handler {
	return &handler{u}
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
