package token

import (
	"time"
)

// Maker is an interface for creating and verifying tokens
type Maker interface {
	// CreateAccess takes an userID and an expiration duration to create a new access token. This
	// short lived token will be used in every request that requires user authentication
	CreateAccess(userID int64, duration time.Duration) (string, error)

	// CreateRefresh takes an userID and expiration to create a longer lived token. This token is used
	// to request a new access token for a valid user
	CreateRefresh(userID int64, duration time.Duration) (string, error)

	// VerifyToken checks if a provided token is valid or not
	VerifyToken(token string) (*Payload, error)
}
