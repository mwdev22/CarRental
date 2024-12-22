package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Addr        string
	DatabaseURI string
	SecretKey   []byte
}

func New() *config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &config{
		Addr:        os.Getenv("ADDR"),
		DatabaseURI: os.Getenv("DATABASE_URI"),
		SecretKey:   []byte(os.Getenv("SECRET_KEY")),
	}
}
