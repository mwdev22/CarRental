package main

import (
	"log"

	"github.com/mwdev22/CarRental/internal/api"
	"github.com/mwdev22/CarRental/internal/config"
	"github.com/mwdev22/CarRental/internal/database"
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
