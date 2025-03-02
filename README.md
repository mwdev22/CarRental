# Car Rental API

## Overview

The **Car Rental API** is a backend service built using **Golang** for managing car rental operations. It includes functionalities such as booking cars, user authentication, and vehicle management. The API is documented using **Swagger** and uses **PostgreSQL** as its database.

## Features

- User authentication and authorization
- Car inventory management
- Rental booking system
- PostgreSQL database with migrations
- RESTful API documented with Swagger

## Prerequisites

Before running the project, ensure you have the following installed:

- **Golang** (>=1.18)
- **PostgreSQL**
- **Make** (for running commands)
- **Golang-Migrate** (database migration tool)
- **Swag** (Swagger documentation generator)

## Setup

### 1. Clone the Repository

```sh
git clone https://github.com/mwdev22/car_rental.git
cd car_rental
```

### 2. Set Up Environment Variables

Create a `.env` file in the project root and add the required variables:

```env
DB_URI=postgres://user:password@localhost:5432/car_rental?sslmode=disable
GOOS=linux
GOARCH=amd64
SECRET_KEY = "4325tcwergtasfsGF453VYRE43YQ34"
ADDR = :8080
```

### 3. Install Dependencies

```sh
go mod tidy
```

### 4. Run Database Migrations

To create new migrations:

```sh
make migrate-create
```

To apply migrations:

```sh
make migrate-up
```

To revert the last migration:

```sh
make migrate-down
```

To drop all migrations:

```sh
make migrate-drop
```

## Running the Application

### Build the Application

```sh
make build
```

### Run the Application

```sh
make run
```

### Run Tests

```sh
make test
```

### Benchmark Tests

```sh
make benchmark
```

## API Documentation

Swagger documentation can be generated using:

```sh
make swag
```

Once generated, access the API documentation at:

```
http://localhost:8080/swagger/index.html
```

## Project Structure

```
├── cmd/                # Main application entry point
├── internal/           # Internal  logic
├── migrations/         # Database migrations
├── .env                # Environment variables
├── Makefile            # Build and migration commands
├── go.mod              # Go dependencies
├── README.md           # Project documentation
```

## Contribution

1. Fork the repository.
2. Create a new branch (`feature-branch-name`).
3. Commit your changes (`git commit -m 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch-name`).
5. Open a pull request.

## License

This project is licensed under the MIT License.
