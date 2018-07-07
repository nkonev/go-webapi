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

			"2_alter.up.sql":
`ALTER TABLE users ADD COLUMN password text;`,

			"3_update.up.sql":
`UPDATE users SET password = 'password'; ALTER TABLE users ALTER COLUMN  password  SET NOT NULL;`,

			"4_update.up.sql":
`UPDATE users SET password = '$2a$10$yZV9IDfxDQjGm1eAvbWip.9JxQzWKTcKm26PGEa/IiMtVc6TJMACu';`,
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