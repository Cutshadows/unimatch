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

func MigrateFs(db *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil) // Reset the base FS after migration
	}()

	return Migrate(db, dir)

}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}
	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
