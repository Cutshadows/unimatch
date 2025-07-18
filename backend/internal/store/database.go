package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib" // Import pgx driver for database/sql
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=centra-db-admin password=wR5BWJGxWbHzse9ltF6kbI0NKPchYVTr dbname=unimatch_db sslmode=disable")

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	fmt.Println("Database connection established successfully %w", db.Ping())
	return db, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}
