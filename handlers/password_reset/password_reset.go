package password_reset

import "github.com/labstack/echo"

type handler struct {

}

func NewHandler() *handler{
	return &handler{}
}

func (h *handler) RequestPasswordReset(c echo.Context) error {
	return nil
}

func (h *handler) ConfirmPasswordReset(c echo.Context) error {
	return nil
}