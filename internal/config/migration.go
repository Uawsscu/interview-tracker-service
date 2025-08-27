package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func migrationDB() (*sql.DB, error) {
	// Create a database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	fmt.Println("[database] Connected to PostgreSQL successfully!")

	// Ping the database to ensure a successful connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	// Specify the path to your migration files (relative to your project directory)
	migrationPath := "file://migrations"

	// Create a new migration instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath, // Path to your migration files
		"postgres",    // Database driver name
		driver,
	)
	if err != nil {
		return nil, err
	}

	// Run the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	// Optionally, you can log the migration status
	version, dirty, err := m.Version()
	if err != nil {
		return nil, err
	}

	log.Printf("Applied migration version: %v (Dirty: %v)\n", version, dirty)

	return db, nil
}
