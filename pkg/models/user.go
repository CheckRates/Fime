package models

import "time"

type User struct {
	ID                int64
	Name              string
	Email             string
	HashedPassword    string
	PasswordChangedAt time.Time
	CreatedAt         time.Time `db:"createdAt"`
}

// Public facing user struct
type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewUserResponse(dbUser User) UserResponse {
	return UserResponse{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
	}
}

// Parameters to list users from a repository
type ListUserParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

// Paramaters to create a new user in the repository
type CreateUserParams struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
}

// Parameters to change a user's name in the repository
type UpdateUserParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
