package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/go-echo-api-test-sample/auth"
	"github.com/go-echo-api-test-sample/db"
	"github.com/go-echo-api-test-sample/handlers/facebook"
	"github.com/go-echo-api-test-sample/handlers/users"
	"github.com/go-echo-api-test-sample/models/session"
	"github.com/go-echo-api-test-sample/models/user"
	"github.com/go-echo-api-test-sample/services"
	"github.com/gobuffalo/packr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

func configureEcho(mailer services.Mailer, facebookClient facebook.FacebookClient, redis *redis.Client,
	sessionModel session.SessionModel, sqlConnection echoConnection) *echo.Echo {


	fromAddress := viper.GetString("mail.registration.fromAddress")
	subject := viper.GetString("mail.registration.subject")
	bodyTemplate := viper.GetString("mail.registration.body.template")
	smtpHostPort := viper.GetString("mail.smtp.address")
	smtpUserName := viper.GetString("mail.smtp.username")
	smtpPassword := viper.GetString("mail.smtp.password")

	confirmationTokenTtl := viper.GetDuration("confirmation.token.ttl")
	sessionTtl := viper.GetDuration("session.ttl")

	url := viper.GetString("url")
	facebookClientId := viper.GetString("facebook.clientId")
	facebookSecret := viper.GetString("facebook.clientSecret")

	db1 := sqlConnection.Connection

	userModel := user.NewUserModel(db1)
	usersHandler := users.NewHandler(userModel)
	fbCallback := "/auth/fb/callback"
	facebookHandler := facebook.NewHandler(facebookClient, facebookClientId, facebookSecret, url+fbCallback, userModel)

	log.SetOutput(os.Stdout)

	static := packr.NewBox("./static")

	e := echo.New()

	e.Use(getAuthMiddleware(sessionModel, stringsToRegexpArray("/user.*", "/auth/.*", "/static.*", "/confirm.*")))
	//e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	e.POST("/auth/login", getLogin(sessionModel, userModel, sessionTtl))
	e.GET("/users/:id", usersHandler.GetDetail)
	e.GET("/users", usersHandler.GetIndex)
	e.GET("/profile", usersHandler.GetProfile)
	e.POST("/auth/register", usersHandler.Register(mailer, fromAddress, subject, bodyTemplate, smtpHostPort, smtpUserName, smtpPassword, url, redis, confirmationTokenTtl))
	e.GET("/confirm/registration", usersHandler.ConfirmRegistration(db1, redis))

	// facebook
	e.Any("/auth/fb", facebookHandler.RedirectForLogin())
	e.Any(fbCallback, facebookHandler.CallBackHandler())

	e.Pre(getStaticMiddleware(static))

	return e
}

func initViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	// path to look for the config file in
	viper.AddConfigPath("./config-dev")
	// call multiple times to add many search paths
	viper.SetEnvPrefix("GO_EXAMPLE")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	// Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
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

func stringsToRegexpArray(strings ...string) []regexp.Regexp {
	regexps := make([]regexp.Regexp, len(strings))
	for i, str := range strings {
		r, err := regexp.Compile(str)
		if err != nil {
			panic(err)
		} else {
			regexps[i] = *r
		}
	}
	return regexps
}

func getLogin(sessionModel session.SessionModel, userModel *user.UserModelImpl, sessionTtl time.Duration) echo.HandlerFunc {
	return func(context echo.Context) error {
		return auth.LoginManager(context, sessionModel, userModel, sessionTtl)
	}
}

func getAuthMiddleware(sm session.SessionModel, whitelist []regexp.Regexp) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return auth.CheckSession(c, next, sm, whitelist)
		}
	}
}

func sessionModel(redis *redis.Client) session.SessionModel {
	return session.SessionModel{Redis: *redis}
}

func main() {
	initViper()
	container := dig.New()
	container.Provide(func() services.Mailer {
		return &services.MailerImpl{}
	})
	container.Provide(func() facebook.FacebookClient {
		return &facebook.FacebookClientImpl{}
	})
	container.Provide(db.ConfigureRedis)
	container.Provide(sessionModel)
	container.Provide(configureEcho)

	container.Provide(db.ConnectDb, dig.Name("migrationSqlxConnection"))
	container.Provide(db.ConnectDb, dig.Name("appSqlxConnection"))

	if migrationErr := container.Invoke(runMigration); migrationErr != nil {
		log.Fatalf("Error during invoke migration: %v", migrationErr)
	}

	if echoErr := container.Invoke(runEcho); echoErr != nil {
		log.Fatalf("Error during invoke echo: %v", echoErr)
	}
	log.Infof("Shutting down")
}

type migrationConnection struct {
	dig.In
	Connection *sqlx.DB `name:"migrationSqlxConnection"`
}

type echoConnection struct {
	dig.In
	Connection *sqlx.DB `name:"appSqlxConnection"`
}

func runMigration(p migrationConnection){
	dropObjects := viper.GetBool("postgresql.dropObjects")
	dropObjectsSql := viper.GetString("postgresql.dropObjectsSql")
	db.MigrateX(p.Connection, dropObjects, dropObjectsSql)
}

// rely on viper import and it's configured by
func runEcho(e *echo.Echo) {
	address := viper.GetString("server.address")
	shutdownTimeout := viper.GetDuration("server.shutdown.timeout")

	log.Info("Starting server")
	// Start server in another goroutine
	go func() {
		if err := e.Start(address); err != nil {
			log.Infof("shutting down the server due error %v", err)
		}
	}()

	log.Info("Waiting for interrupt (2) (Ctrl+C)")
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Infof("Got signal %v - will forcibly close after %v", os.Interrupt, shutdownTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel() // releases resources if slowOperation completes before timeout elapses
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}