package main

import (
	"github.com/labstack/echo"
	"github.com/go-echo-api-test-sample/handlers"
	"github.com/go-echo-api-test-sample/models/user"
	"github.com/go-echo-api-test-sample/db"
	"github.com/go-echo-api-test-sample/migrations"
	"os"
	"github.com/labstack/gommon/log"
	"os/signal"
	"context"
	"time"
	"github.com/labstack/echo/middleware"
	_ "github.com/volatiletech/authboss/auth"
	"github.com/go-echo-api-test-sample/auth"
	"github.com/go-echo-api-test-sample/models/session"
)

func configureEcho() *echo.Echo {
	d0 := db.DBConnect()
	migrations.MigrateX(d0)
	d1 := db.DBConnect()
	m := user.NewUserModel(d1)
	h := users.NewHandler(m)

	r := db.Connect()
	sm := session.SessionModel{Redis: *r}

	log.SetOutput(os.Stdout)

	e := echo.New()

	e.Use(getAuthMiddleware(sm, []string{"/user.*", "/auth2/.*"}))
	//e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	e.POST("/auth2/login", getLogin(sm, m))
	e.GET("/users/:id", h.GetDetail)
	e.GET("/users", h.GetIndex)

	return e
}

func getLogin(sessionModel session.SessionModel, userModel *user.UserModelImpl) echo.HandlerFunc {
	return func (context echo.Context) error {
		return auth.LoginManager(context, sessionModel, userModel)
	}
}


func getAuthMiddleware(sm session.SessionModel, whitelist []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return auth.CheckSession(c, next, sm, whitelist);
		}
	}
}

func main() {
	e := configureEcho()

	log.Info("Starting server")
	// Start server
	go func() {
		if err := e.Start(":1234"); err != nil {
			log.Infof("shutting down the server due error %v", err)
		}
	}()

	log.Info("Waiting for interrupt (2) (Ctrl+C)")
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Infof("Got signal %v", os.Interrupt)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
