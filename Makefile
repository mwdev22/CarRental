# Include .env file
include .env
export $(shell sed 's/=.*//' .env)

# Project variables
PROJECT_NAME := file_storage
PKG := ./...
MAIN := ./cmd/
MIGRATIONS_DIR := migrations

# Go commands
BUILD := go build
CLEAN := go clean
FMT := go fmt
VET := go vet
TEST := go test
RUN := go run
MIGRATE := migrate

# Targets
.PHONY: all build clean fmt vet test run migrate-create migrate-up migrate-down migrate-drop

# Default target
all: fmt vet test build

# Build the application
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(BUILD) -o $(PROJECT_NAME).exe $(MAIN)

# Clean the build files
clean:
	$(CLEAN)
	rm -f $(PROJECT_NAME)

# Format the Go code
fmt:
	$(FMT) $(PKG)

# Vet the Go code
vet:
	$(VET) $(PKG)

# Run tests
test:
	$(TEST) $(PKG)

# Run the application
run:
	$(RUN) $(MAIN)

# Migration commands
migrate-create:
	@export DB_SOURCE=$(DB_URI); \
	read -p "Enter migration name: " name; \
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

migrate-up:
	@export DB_SOURCE=$(DB_URI); \
	$(MIGRATE) -database $$DB_SOURCE -path $(MIGRATIONS_DIR) up

migrate-down:
	@export DB_SOURCE=$(DB_URI); \
	$(MIGRATE) -database $$DB_SOURCE -path $(MIGRATIONS_DIR) down 1

migrate-force:
	@export DB_SOURCE=$(DB_URI); \
	$(MIGRATE) -database $$DB_SOURCE -path $(MIGRATIONS_DIR) force $(VERSION)

migrate-drop:
	@export DB_SOURCE=$(DB_URI); \
	$(MIGRATE) -database $$DB_SOURCE -path $(MIGRATIONS_DIR) drop

swag:
	@swag init --generalInfo ./cmd/main.go
