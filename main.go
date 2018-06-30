package main

import (
	"github.com/labstack/echo"
	"github.com/go-echo-api-test-sample/handlers"
	"github.com/go-echo-api-test-sample/models"
	"github.com/go-echo-api-test-sample/db"
	"github.com/go-echo-api-test-sample/migrations"
	"os"
	"github.com/labstack/gommon/log"
	"os/signal"
	"context"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)

	e := echo.New()

	d0 := db.DBConnect()
	migrations.MigrateX(d0)
	d1 := db.DBConnect()
	h := users.NewHandler(user.NewUserModel(d1))

	e.GET("/users", h.GetIndex)
	e.GET("/users/:id", h.GetDetail)

	//e.Logger.Fatal(e.Start(":1234"))

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil {
			// TODO not works
			// https://echo.labstack.com/cookbook/graceful-shutdown
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
