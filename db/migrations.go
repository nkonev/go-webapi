package db

import (
	"database/sql"
	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/gobuffalo/packr"
)

func _migrate(db *sql.DB, dropObjects bool, dropObjectsSql string) {
	if dropObjects {
		if _, err := db.Exec(dropObjectsSql); err != nil {
			log.Errorf("Error during dropping tables", err)
		}
	}

	driver, err0 := postgres.NewFromDB(db)
	defer driver.Close()

	if err0 != nil {
		log.Fatalf("Unable to open connection to postgres server: %s", err0)
	}

	packrSource := &migration.PackrMigrationSource{
		Box: packr.NewBox("migrations"),
	}

	applied, err := migration.Migrate(driver, packrSource, migration.Up, 0)
	if err != nil {
		log.Fatalf("Error during migration %s", err)
	}

	log.Infof("Applied %d migrations", applied)
}

func MigrateX(db *sqlx.DB, dropObjects bool, dropObjectsSql string){
	_migrate(db.DB, dropObjects,dropObjectsSql)
}