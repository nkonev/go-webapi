package db

import (
	"github.com/labstack/gommon/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDb(postgresqlConnectString string) *sqlx.DB {
	// connect_timeout ./vendor/github.com/lib/pq/doc.go in seconds
	// statement_timeout https://postgrespro.ru/docs/postgrespro/9.6/runtime-config-client
	db, err := sqlx.Connect("postgres", postgresqlConnectString)
	if err != nil {
		log.Panic(err)
	}
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(500)
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	return db
}
