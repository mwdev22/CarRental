package store

import (
	"context"
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

func (r *UserRepo) Create(ctx context.Context, u *User) error {
	query := `INSERT INTO users (username, password, email, role) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, u.Username, u.Password, u.Email, u.Role)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	query := `SELECT id, username, password, email, role FROM users WHERE username = $1`
	err := r.DB.Get(&user, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %v", err)
	}
	return &user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*User, error) {
	var user User
	query := `SELECT id, username, email, role FROM users WHERE id = $1`
	err := r.DB.Get(&user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %v", err)
	}
	return &user, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("user with provided id not found")
	}
	return nil
}

func (r *UserRepo) Update(ctx context.Context, u *User) error {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3`
	_, err := r.DB.Exec(query, u.Username, u.Email, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}
