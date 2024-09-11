package database

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

func RunMigrations(dsn string) {
	db, err := sql.Open(
		"postgres",
		dsn,
	)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(
		db,
		&postgres.Config{},
	)
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://service/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("Migrations applied successfully!")
}
