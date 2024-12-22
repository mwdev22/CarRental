package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) CreateUser(u *User) error {
	query := `INSERT INTO users (username, password, email, created) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, u.Username, u.Password, u.Email, u.Created)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (r *UserRepo) GetByUsername(username string) (*User, error) {
	var user User
	query := `SELECT id, username, password, email FROM users WHERE username = $1`
	err := r.DB.Get(&user, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %v", err)
	}
	return &user, nil
}
