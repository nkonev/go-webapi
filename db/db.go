package db

import (
	"github.com/labstack/gommon/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func ConnectDbInternal(postgresqlConnectString string, maxPostgreConns int, minPostgreConns int) *sqlx.DB {
	// connect_timeout ./vendor/github.com/lib/pq/doc.go in seconds
	// statement_timeout https://postgrespro.ru/docs/postgrespro/9.6/runtime-config-client
	db, err := sqlx.Connect("postgres", postgresqlConnectString)
	if err != nil {
		log.Panic(err)
	}
	db.SetMaxOpenConns(maxPostgreConns)
	db.SetMaxIdleConns(minPostgreConns)
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	return db
}

func ConnectDb() *sqlx.DB {
	postgresqlConnectString := viper.GetString("postgresql.connectString")
	maxPostgreConns := viper.GetInt("postgresql.maxOpenConnections")
	minPostgreConns := viper.GetInt("postgresql.minOpenConnections")

	return ConnectDbInternal(postgresqlConnectString, maxPostgreConns, minPostgreConns)
}