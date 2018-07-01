package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBConnect() *sqlx.DB {
	// connect_timeout ./vendor/github.com/lib/pq/doc.go in seconds
	// statement_timeout https://postgrespro.ru/docs/postgrespro/9.6/runtime-config-client
	db, err := sqlx.Connect("postgres", "host=172.24.0.2 user=postgres password=postgresqlPassword dbname=postgres connect_timeout=2 statement_timeout=2000 sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	return db
}
