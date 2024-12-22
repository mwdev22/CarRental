package main

import (
	"log"

	"github.com/mwdev22/FileStorage/internal/api"
	"github.com/mwdev22/FileStorage/internal/config"
	"github.com/mwdev22/FileStorage/internal/database"
)

func main() {
	cfg := config.New()
	db, err := database.OpenSQLConnection(cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	api := api.New(cfg.Addr, db)

	if err := api.Start(); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}
