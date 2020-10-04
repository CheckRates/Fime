package postgres

import (
	"github.com/checkrates/Fime/fime"
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
		return fime.User{}, err
	}
	return u, nil
}

//Users retrieve all users
func (s *UserStore) Users(limit int, offset int) ([]fime.User, error) {
	var uu []fime.User
	if err := s.Select(&uu, `SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return []fime.User{}, err
	}
	return uu, nil
}

//CreateUser creates a user in the database
func (s *UserStore) CreateUser(u *fime.User) error {
	err := s.Get(u, `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING *`,
		u.Name, u.Email, u.Password)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser updates info about a existing user
func (s *UserStore) UpdateUser(u *fime.User) error {
	if err := s.Get(u, `UPDATE users SET name=$1 WHERE id=$2 RETURNING *`, u.Name, u.ID); err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user from the database
func (s *UserStore) DeleteUser(id int64) error {
	if _, err := s.Exec(`DELETE FROM users WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
