package env

import (
	"database/sql"
	"github.com/labstack/gommon/log"
)

type Env struct {
	db *sql.DB
	logger *log.Logger
}