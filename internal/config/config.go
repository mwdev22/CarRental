package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var SecretKey []byte

type config struct {
	Addr        string
	DatabaseURI string
}

func New() *config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	SecretKey = []byte(os.Getenv("SECRET_KEY"))

	return &config{
		Addr:        os.Getenv("ADDR"),
		DatabaseURI: os.Getenv("DATABASE_URI"),
	}
}
