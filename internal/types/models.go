package types

import (
	"time"
)

type User struct {
	ID       int       `db:"id" json:"id"`
	Username string    `db:"username" json:"username"`
	Password []byte    `db:"password" json:"-"`
	Email    string    `db:"email" json:"email"`
	Role     UserRole  `db:"role" json:"role"`
	Created  time.Time `db:"created" json:"-"`
}

type Company struct {
	ID      int       `json:"id" db:"id"`             // Unique ID for the company
	OwnerID int       `json:"owner_id" db:"owner_id"` // ID of the user who owns the company
	Name    string    `json:"name" db:"name"`         // Company name
	Email   string    `json:"email" db:"email"`       // Contact email
	Phone   string    `json:"phone" db:"phone"`       // Contact phone number
	Address string    `json:"address" db:"address"`   // Address of the company
	Created time.Time `json:"created_at" db:"created_at"`
	Updated time.Time `json:"updated_at" db:"updated_at"` // Last updated timestamp
}

type Car struct {
	ID             int       `json:"id" db:"id"`
	CompanyID      int       `json:"company_id" db:"company_id"` // ID of the company that owns the car
	Make           string    `json:"make" db:"make"`             // Make (e.g., Toyota, Ford)
	Model          string    `json:"model" db:"model"`           // Model (e.g., Corolla, Mustang)
	Year           int       `json:"year" db:"year"`             // Year of manufacture
	Color          string    `json:"color" db:"color"`
	RegistrationNo string    `json:"registration_no" db:"registration_no"` // Car registration number
	PricePerDay    float64   `json:"price_per_day" db:"price_per_day"`
	Created        time.Time `json:"created_at" db:"created_at"`
	Updated        time.Time `json:"updated_at" db:"updated_at"` // Last updated timestamp
}

type Booking struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	CarID     int       `json:"car_id" db:"car_id"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	Total     float64   `json:"total" db:"total"`
}
