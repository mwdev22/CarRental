package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register the PostgreSQL driver
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

func OpenTestSqlDB(uri string) (*sqlx.DB, string, error) {
	var db *sqlx.DB
	var err error

	// scrape connection string without db name
	uriParts := strings.Split(uri, "/")
	uriWithoutDB := strings.Join(uriParts[:len(uriParts)-1], "/")

	db, err = sqlx.Open("postgres", fmt.Sprintf("%s/postgres", uriWithoutDB))
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
		return nil, "", err
	}
	// create a new test database
	testDBName := "test_" + time.Now().Format("20060102150405")
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", testDBName))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create test database: %w", err)
	}

	db.Close()

	// connect to the created test database
	db, err = OpenSQLConnection(fmt.Sprintf("%s/%s", uriWithoutDB, testDBName))
	if err != nil {
		log.Fatalf("failed to connect to the test database: %v", err)
		return nil, "", err
	}

	// find the migrations directory
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
		return nil, "", err
	}

	migrationsPath := filepath.Join(basePath, "migrations")

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("%s/%s", uriWithoutDB, testDBName),
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to initialize migrations: %w", err)
	}
	// perform migrations on test database
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, "", fmt.Errorf("failed to apply migrations: %w", err)
	}

	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(30)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, testDBName, nil
}
