package store

import (
	"time"

	"github.com/mwdev22/FileStorage/internal/types"
)

type User struct {
	ID       int            `db:"id" json:"id"`
	Username string         `db:"username" json:"username"`
	Password []byte         `db:"password" json:"-"`
	Email    string         `db:"email" json:"email"`
	Role     types.UserRole `db:"role" json:"role"`
	Created  time.Time      `db:"created" json:"-"`
}

type Company struct {
	ID      int       `json:"id" db:"id"`           // Unique ID for the company
	Name    string    `json:"name" db:"name"`       // Company name
	Email   string    `json:"email" db:"email"`     // Contact email
	Phone   string    `json:"phone" db:"phone"`     // Contact phone number
	Address string    `json:"address" db:"address"` // Address of the company
	Created time.Time `json:"created_at" db:"created_at"`
	Updated time.Time `json:"updated_at" db:"updated_at"` // Last updated timestamp
}

type Car struct {
	ID             int       `json:"id" db:"id"`
	Make           string    `json:"make" db:"make"`   // Make (e.g., Toyota, Ford)
	Model          string    `json:"model" db:"model"` // Model (e.g., Corolla, Mustang)
	Year           int       `json:"year" db:"year"`   // Year of manufacture
	Color          string    `json:"color" db:"color"`
	RegistrationNo string    `json:"registration_no" db:"registration_no"` // Car registration number
	PricePerDay    float64   `json:"price_per_day" db:"price_per_day"`
	Created        time.Time `json:"created_at" db:"created_at"`
	Updated        time.Time `json:"updated_at" db:"updated_at"` // Last updated timestamp
}
