package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var SecretKey []byte

type config struct {
	Addr        string
	DatabaseURI string
	TestDbURI   string
}

func New() *config {
	projectDir, err := findProjectRoot()

	if err != nil {
		log.Fatalf("failed to locate project root: %v", err)
	}
	if err := os.Chdir(projectDir); err != nil {
		log.Fatalf("failed to set working directory: %v", err)
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	SecretKey = []byte(os.Getenv("SECRET_KEY"))

	return &config{
		Addr:        os.Getenv("ADDR"),
		DatabaseURI: os.Getenv("DB_URI"),
	}
}

// helper function to find the project root (usefull for tests to locate .env file)
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(dir + "/go.mod"); err == nil {
			return dir, nil
		}
		if dir == "/" {
			return "", os.ErrNotExist
		}
		// if go.mod doesnt exist in current directory, move up one
		dir = filepath.Dir(dir)
	}
}
