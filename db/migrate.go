package db

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func RunMigration(migrationURL *string, dbSource *string) {
	migration, err := migrate.New(*migrationURL, *dbSource)
	if err != nil {
		log.Fatal("cannot create migration: ", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("cannot run migration: ", err)
	}
	log.Println("migration complete")
}
