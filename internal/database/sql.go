package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register the PostgreSQL driver
	_ "github.com/mattn/go-sqlite3"
)

func OpenSQLConnection(uri string) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	db, err = sqlx.Open("postgres", uri)

	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping the database: %v", err)
		return nil, err
	}

	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(30)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil

}

// test db is an sqlite database
func OpenTestSqlDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatalf("failed to connect to SQLite database: %v", err)
		return nil, err
	}

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping SQLite database: %v", err)
		return nil, err
	}
	db.MustExec("PRAGMA foreign_keys = ON;")
	migrationsPath := filepath.Join(basePath, "test_migrations")
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("sqlite3://%s/test.db", basePath),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrations: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to apply SQLite migrations: %w", err)
	}
	log.Println("Migrations applied successfully")

	var tables []string
	err = db.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		log.Printf("failed to list tables in SQLite database: %v", err)
	} else {
		log.Printf("Tables in database: %v", tables)
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
