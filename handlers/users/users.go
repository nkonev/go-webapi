package users

import (
	"github.com/labstack/echo"
	"github.com/nkonev/go-webapi/models/user"
	"github.com/nkonev/go-webapi/utils"
	"net/http"
	"strconv"
)

type resultLists struct {
	Users []user.User `json:"users"`
}

type userHandler struct {
	UserModel user.UserModel
}

func NewUserHandler(u user.UserModel) *userHandler {
	return &userHandler{u}
}

type RegisterDTO struct {
	Email    string // email
	Password string
}

func (h *userHandler) GetIndex(c echo.Context) error {
	lists, e := h.UserModel.FindAll()
	if e != nil {
		return c.JSON(http.StatusInternalServerError, utils.H{"error": "We had error ("})
	}
	u := &resultLists{
		Users: lists,
	}
	return c.JSON(http.StatusOK, u)
}

func (h *userHandler) GetDetail(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	u, e := h.UserModel.FindByID(idInt)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, utils.H{"error": "We had error ("})
	}
	return c.JSON(http.StatusOK, u)
}

func (h *userHandler) GetProfile(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.H{"message": "You see your profile"})
}

