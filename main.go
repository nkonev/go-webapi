package main

import (
	"github.com/labstack/echo"
	"github.com/go-echo-api-test-sample/handlers"
	"github.com/go-echo-api-test-sample/models/user"
	"github.com/go-echo-api-test-sample/db"
	"os"
	"github.com/labstack/gommon/log"
	"os/signal"
	"context"
	"time"
	"github.com/labstack/echo/middleware"
	"github.com/go-echo-api-test-sample/auth"
	"github.com/go-echo-api-test-sample/models/session"
	"github.com/spf13/viper"
	"fmt"
	"github.com/go-echo-api-test-sample/services"
	"regexp"
	"github.com/gobuffalo/packr"
	"net/http"
	"strings"
)

func configureEcho(mailer services.Mailer) *echo.Echo {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")   // path to look for the config file in
	viper.AddConfigPath("./config-dev")  // call multiple times to add many search paths
	viper.SetEnvPrefix("GO_EXAMPLE")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	redisAddr := viper.GetString("redis.addr")
	redisPassword := viper.GetString("redis.password")
	redisDbNum := viper.GetInt("redis.db")
	redisFlushOnStart := viper.GetBool("redis.flushOnStart")

	postgresqlConnectString := viper.GetString("postgresql.connectString")
	maxPostgreConns := viper.GetInt("postgresql.maxOpenConnections")
	minPostgreConns := viper.GetInt("postgresql.minOpenConnections")
	dropObjects := viper.GetBool("postgresql.dropObjects")
	dropObjectsSql := viper.GetString("postgresql.dropObjectsSql")

	fromAddress := viper.GetString("mail.registration.fromAddress")
	subject := viper.GetString("mail.registration.subject")
	bodyTemplate := viper.GetString("mail.registration.body.template")
	smtpHostPort := viper.GetString("mail.smtp.address")
	smtpUserName := viper.GetString("mail.smtp.username")
	smtpPassword := viper.GetString("mail.smtp.password")

	url := viper.GetString("url")

	d0 := db.ConnectDb(postgresqlConnectString, maxPostgreConns, minPostgreConns)
	db.MigrateX(d0, dropObjects, dropObjectsSql)
	d1 := db.ConnectDb(postgresqlConnectString, maxPostgreConns, minPostgreConns)
	m := user.NewUserModel(d1)
	h := users.NewHandler(m)

	r := db.ConnectRedis(redisAddr, redisPassword, redisDbNum, redisFlushOnStart)
	sm := session.SessionModel{Redis: *r}

	log.SetOutput(os.Stdout)

	static := packr.NewBox("./static")

	e := echo.New()

	e.Use(getAuthMiddleware(sm, stringsToRegexpArray("/user.*", "/auth2/.*", "/static.*")))
	//e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	e.POST("/auth2/login", getLogin(sm, m))
	e.GET("/users/:id", h.GetDetail)
	e.GET("/users", h.GetIndex)
	e.GET("/profile", h.GetProfile)
	e.POST("/auth2/register", h.Register(mailer, fromAddress, subject, bodyTemplate, smtpHostPort, smtpUserName, smtpPassword, url, r))

	e.Pre(getStaticMiddleware(static))

	return e
}

func getStaticMiddleware(box packr.Box) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqUrl := c.Request().RequestURI
			if reqUrl == "/" || reqUrl == "/index.html" || strings.HasPrefix(reqUrl, "/assets") {
				http.FileServer(box).
					ServeHTTP(c.Response().Writer, c.Request())
				return nil
			} else {
				return next(c)
			}
		}
	}
}

func stringsToRegexpArray(strings... string) []regexp.Regexp {
	regexps := make([]regexp.Regexp, len(strings))
	for i, str := range strings {
		r, err := regexp.Compile(str)
		if err != nil {
			panic(err)
		} else {
			regexps[i] = *r;
		}
	}
	return regexps
}

func getLogin(sessionModel session.SessionModel, userModel *user.UserModelImpl) echo.HandlerFunc {
	return func (context echo.Context) error {
		return auth.LoginManager(context, sessionModel, userModel)
	}
}


func getAuthMiddleware(sm session.SessionModel, whitelist []regexp.Regexp) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return auth.CheckSession(c, next, sm, whitelist);
		}
	}
}

func main() {
	e := configureEcho(&services.MailerImpl{})

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
