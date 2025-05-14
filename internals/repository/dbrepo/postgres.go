package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/kiniconnect/booking-app/internals/models"
)

// AllUsers returns true if all users are returned
func (m *PostgresDBRepo) AllUsers() bool {
	// check if the database connection is nil
	if m.DB == nil {
		return false
	}
	fmt.Println("DB connection is not nil")

	// query to get all users
	query := `SELECT id, first_name, last_name, email, password, created_at, updated_at FROM users`

	rows, err := m.DB.Query(query)
	if err != nil {
		return false
	}
	defer rows.Close()

	// loop through the rows and scan the values into a User struct
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return false
		}
	}

	// check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return false
	}

	fmt.Println("All users returned successfully")

	return true
}

// InsertReservation insert a reservation to the database
func (m *PostgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var NewID int

	stmt := ` insert into reservations (firstname, lastname, phone, email, start_date, end_date, room_id, created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Phone,
		res.Email,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&res.ID)

	if err != nil {
		return 0, err
	}

	NewID = res.ID
	// check if the reservation ID is valid
	if NewID <= 0 {
		return 0, fmt.Errorf("invalid reservation ID: %d", NewID)
	}
	return NewID, nil
}

func (m *PostgresDBRepo) InsertRoomRestriction(roomRestriction models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `INSERT INTO room_restrictions (start_date, end_date, room_id, 
				reservation_id, created_at, updated_at, restriction_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		roomRestriction.StartDate,
		roomRestriction.EndDate,
		roomRestriction.RoomID,
		roomRestriction.ReservationID,
		time.Now(),
		time.Now(),
		roomRestriction.RestrictionID,
	)
	if err != nil {
		return fmt.Errorf("error inserting room restriction: %w", err)
	}

	return nil
}
