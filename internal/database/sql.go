package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register the PostgreSQL driver
)

func OpenSQLConnection(uri string) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	db, err = sqlx.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

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

	fmt.Println("Connected to database...\n")
	return db, nil

}
