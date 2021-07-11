package postgres

import (
	"github.com/jmoiron/sqlx"
)

// NewUserStore returns the access point to all the users of Fime
func NewUserStore(db *sqlx.DB) *UserSQL {
	return &UserSQL{
		DB: db,
	}
}

// UserSQL is the database access point to users
type UserSQL struct {
	*sqlx.DB
}

// CreateUserParams provides all info to create a new user in the db
type CreateUserParams struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
}

// UpdateUserParams provides all info to change a user's name in the db
type UpdateUserParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// User retrieves a user from the database by id
func (s *UserSQL) User(id int64) (User, error) {
	var u User
	if err := s.Get(&u, `SELECT * FROM users WHERE id=$1 LIMIT 1`, id); err != nil {
		return User{}, err
	}
	return u, nil
}

// UserByEmail retrieves a user from the database by email
func (s *UserSQL) UserByEmail(email string) (User, error) {
	var u User
	if err := s.Get(&u, `SELECT * FROM users WHERE email=$1 LIMIT 1`, email); err != nil {
		return User{}, err
	}
	return u, nil
}

//Users retrieve all users
func (s *UserSQL) Users(args ListParams) ([]User, error) {
	uu := []User{}
	if err := s.Select(&uu, `SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2`, args.Limit, args.Offset); err != nil {
		return []User{}, err
	}
	return uu, nil
}

// CreateUser creates a user in the database
func (s *UserSQL) CreateUser(args CreateUserParams) (User, error) {
	var u User
	err := s.Get(&u, `INSERT INTO users (name, email, hashedPassword) VALUES ($1, $2, $3) RETURNING *`,
		args.Name, args.Email, args.HashedPassword)
	if err != nil {
		return u, err
	}
	return u, nil
}

// UpdateUser updates info about a existing user
func (s *UserSQL) UpdateUser(args UpdateUserParams) (User, error) {
	var u User
	if err := s.Get(&u, `UPDATE users SET name=$1 WHERE id=$2 RETURNING *`, args.Name, args.ID); err != nil {
		return u, err
	}
	return u, nil
}

// DeleteUser deletes a user from the database
func (s *UserSQL) DeleteUser(id int64) error {
	if _, err := s.Exec(`DELETE FROM users WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
