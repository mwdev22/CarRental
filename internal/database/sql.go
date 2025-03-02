package database

import (
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func OpenSQLConnection(uri string, dbType string) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	fmt.Println(uri, dbType)

	db, err = sqlx.Open(dbType, uri)

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
