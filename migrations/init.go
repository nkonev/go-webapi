package migrations

import (
	"github.com/pressly/goose"
	"log"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"os"
)

func Migrate(db *sql.DB) {
	dir := "./migrations";
	goose.SetLogger(log.New(os.Stdout, "", log.LstdFlags))
	if err := goose.Run("up", db, dir); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}

func MigrateX(db *sqlx.DB){
	Migrate(db.DB)
}