package migrations

import (
	"database/sql"
	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

func _migrate(db *sql.DB) {
	driver, err0 := postgres.NewFromDB(db)
	defer driver.Close()

	if err0 != nil {
		log.Fatalf("Unable to open connection to postgres server: %s", err0)
	}

	migrations2 := &migration.MemoryMigrationSource{
		Files: map[string]string{
			"1_init.up.sql":

`CREATE TABLE users (
    id int NOT NULL PRIMARY KEY,
    name text,
    surname text,
	lastname text
);
INSERT INTO users VALUES
(0, 'root', 's', 'l'),
(1, 'vojtechvitek', '', '');`,

		},
	}

	applied, err := migration.Migrate(driver, migrations2, migration.Up, 0)
	if err != nil {
		log.Fatalf("Error during migration %s", err)
	}

	log.Infof("Applied %d migrations", applied)
}

func MigrateX(db *sqlx.DB){
	_migrate(db.DB)
}