package main

import (
	"github.com/labstack/echo"
	"github.com/go-echo-api-test-sample/handlers"
	"github.com/go-echo-api-test-sample/models"
	"github.com/go-echo-api-test-sample/db"
	"github.com/go-echo-api-test-sample/migrations"
)

func main() {
	e := echo.New()

	d := db.DBConnect()
	migrations.MigrateX(d)
	h := users.NewHandler(user.NewUserModel(d))

	e.GET("/users", h.GetIndex)
	e.GET("/users/:id", h.GetDetail)

	e.Logger.Fatal(e.Start(":1324"))
}
