package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mwdev22/CarRental/internal/types"
)

type BookingRepositorySQL struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) *BookingRepositorySQL {
	return &BookingRepositorySQL{
		db: db,
	}
}

func (bs *BookingRepositorySQL) Create(ctx context.Context, booking *types.Booking) error {
	query := `INSERT INTO booking (user_id, car_id, start_date, end_date, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	_, err := bs.db.Exec(query, booking.UserID, booking.CarID, booking.StartDate, booking.EndDate, booking.Total)

	if err != nil {
		return fmt.Errorf("error creating booking: %w", err)
	}
	return nil
}

func (bs *BookingRepositorySQL) GetByID(ctx context.Context, id int) (*types.Booking, error) {
	query := `SELECT id, user_id, car_id, start_date, end_date, total  FROM booking WHERE id = $1`
	var booking types.Booking
	err := bs.db.Get(&booking, query, id)
	if err != nil {
		return nil, fmt.Errorf("error getting booking: %w", err)
	}
	return &booking, nil
}

func (bs *BookingRepositorySQL) Update(ctx context.Context, booking *types.Booking) error {
	query := `UPDATE booking SET user_id=$1, car_id=$2, start_date=$3, end_date=$4, total=$5 WHERE id=$6`
	_, err := bs.db.Exec(query, booking.UserID, booking.CarID, booking.StartDate, booking.EndDate, booking.Total, booking.ID)
	if err != nil {
		return fmt.Errorf("error updating bookings: %w", err)
	}
	return nil
}

func (bs *BookingRepositorySQL) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM booking WHERE id=$1`
	_, err := bs.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting bookings: %w", err)
	}
	return nil
}

func (bs *BookingRepositorySQL) GetByUserID(ctx context.Context, userID int) ([]*types.Booking, error) {
	query := `SELECT id, user_id, car_id, start_date, end_date, total FROM booking WHERE user_id = $1`
	var booking []*types.Booking
	err := bs.db.Select(&booking, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting bookings: %w", err)
	}
	return booking, nil
}

func (bs *BookingRepositorySQL) CheckDateAvailability(ctx context.Context, carID int, startDate, endDate time.Time) bool {
	query := `SELECT id FROM booking WHERE car_id = $1 AND start_date <= $2 AND end_date >= $3`
	var booking types.Booking
	err := bs.db.Get(&booking, query, carID, startDate, endDate)
	return err != nil
}

func (bs *BookingRepositorySQL) GetCurrent(ctx context.Context) ([]*types.Booking, error) {
	query := `SELECT id, user_id, car_id, start_date, end_date, total FROM booking WHERE start_date <= $1 AND end_date >= $2`
	var booking []*types.Booking
	err := bs.db.Select(&booking, query, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("error getting bookings: %w", err)
	}
	return booking, nil
}
