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
	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss/defaults"
	_ "github.com/volatiletech/authboss/auth"
	"net/http"
	"github.com/go-echo-api-test-sample/auth"
	"github.com/go-echo-api-test-sample/models/session"
)

func configureEcho() *echo.Echo {
	d0 := db.DBConnect()
	migrations.MigrateX(d0)
	d1 := db.DBConnect()
	m := user.NewUserModel(d1)

	r := db.Connect()
	sm := session.SessionModel{Redis: *r}

	log.SetOutput(os.Stdout)

	authPathPrefix := "/auth2"

	ab := authboss.New()
	ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}
	ab.Config.Storage.Server = &auth.MyServerStorer{Model:*m}
	ab.Config.Storage.SessionState = &auth.MySessionStorer{Model: sm}
	//ab.Config.Storage.CookieState = myCookieImplementation //todo implement
	defaults.SetCore(&ab.Config, true, true)
	if err := ab.Init("auth"); err != nil {
		log.Panic(err)
	}

	e := echo.New()

	//e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	h := users.NewHandler(user.NewUserModel(d1))

	e.GET("/users", h.GetIndex)
	e.GET("/users/:id", h.GetDetail)

	g := e.Group(authPathPrefix)
	g.Use(wrap(http.StripPrefix(authPathPrefix, ab.Config.Core.Router), ab))

	return e
}

func wrap(h http.Handler, ab *authboss.Authboss) echo.MiddlewareFunc {
	h2 := ab.LoadClientStateMiddleware(h)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h2.ServeHTTP(c.Response(), c.Request())
			return next(c)
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
