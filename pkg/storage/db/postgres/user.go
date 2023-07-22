package postgres

import (
	"github.com/checkrates/Fime/pkg/models"
	"github.com/jmoiron/sqlx"
)

type UserSQL struct {
	*sqlx.DB
}

// Returns the access point to Fime's users
func NewUserRepository(db *sqlx.DB) *UserSQL {
	return &UserSQL{
		DB: db,
	}
}

// Retrieves a user from the database by id, if found
func (s *UserSQL) FindById(id int64) (models.User, error) {
	var u models.User
	if err := s.Get(&u, `SELECT * FROM users WHERE id=$1 LIMIT 1`, id); err != nil {
		return models.User{}, err
	}
	return u, nil
}

// Retrieves a user from the database by email, if found
func (s *UserSQL) FindByEmail(email string) (models.User, error) {
	var u models.User
	if err := s.Get(&u, `SELECT * FROM users WHERE email=$1 LIMIT 1`, email); err != nil {
		return models.User{}, err
	}
	return u, nil
}

// Retrieve a subset of users
func (s *UserSQL) GetMultiple(args models.ListUserParams) ([]models.User, error) {
	uu := []models.User{}
	if err := s.Select(&uu, `SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2`, args.Limit, args.Offset); err != nil {
		return []models.User{}, err
	}
	return uu, nil
}

// Creates a user in the database
func (s *UserSQL) Create(args models.CreateUserParams) (models.User, error) {
	var u models.User
	err := s.Get(&u, `INSERT INTO users (name, email, hashedPassword) VALUES ($1, $2, $3) RETURNING *`,
		args.Name, args.Email, args.HashedPassword)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Updates info about a existing user
func (s *UserSQL) Update(args models.UpdateUserParams) (models.User, error) {
	var u models.User
	if err := s.Get(&u, `UPDATE users SET name=$1 WHERE id=$2 RETURNING *`, args.Name, args.ID); err != nil {
		return u, err
	}
	return u, nil
}

// Deletes a user from the database
func (s *UserSQL) Delete(id int64) error {
	if _, err := s.Exec(`DELETE FROM users WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
