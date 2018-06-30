package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func DBConnect() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "host=172.24.0.2 user=postgres password=postgresqlPassword dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

