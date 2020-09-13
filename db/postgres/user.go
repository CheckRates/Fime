package postgres

import (
	"fmt"

	"github.com/checkrates/Fime/fime"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// NewUserStore returns the access point to all the users of Fime
func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

// UserStore is the database access point to users
type UserStore struct {
	*sqlx.DB
}

// User retrieves a user from the database by id
func (s *UserStore) User(id int64) (fime.User, error) {
	var u fime.User
	if err := s.Get(&u, `SELECT * FROM users WHERE id=$1 LIMIT 1`, id); err != nil {
		return fime.User{}, fmt.Errorf("error retrieving user: %w", err)
	}
	return u, nil
}

//Users retrieve all users
func (s *UserStore) Users(limit int, offset int) ([]fime.User, error) {
	var uu []fime.User
	if err := s.Get(&uu, `SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return []fime.User{}, fmt.Errorf("error retrieving users: %w", err)
	}
	return uu, nil
}

//CreateUser creates a user in the database
func (s *UserStore) CreateUser(u *fime.User) error {
	if err := s.Get(u, `INSERT INTO users (name) VALUES ($1) RETURNING *`, u.Name); err != nil {
		return fmt.Errorf("error inserting new user: %w", err)
	}
	return nil
}

// UpdateUser updates info about a existing user
func (s *UserStore) UpdateUser(u *fime.User) error {
	if err := s.Get(u, `UPDATE users SET name = $1 RETURNING *`, u.Name); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user from the database
func (s *UserStore) DeleteUser(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM users WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}