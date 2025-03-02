package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var SecretKey []byte

type config struct {
	Addr        string
	DatabaseURI string
	DBType      string
}

func New() *config {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUri := os.Getenv("DB_URI")

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = strings.Split(dbUri, ":")[0]
	}

	SecretKey = []byte(os.Getenv("SECRET_KEY"))

	return &config{
		Addr:        os.Getenv("ADDR"),
		DatabaseURI: dbUri,
		DBType:      dbType,
	}
}
