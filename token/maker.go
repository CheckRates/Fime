package token

import (
	"time"
)

// Maker is an interface for creating and verifying tokens
type Maker interface {
	// CreateToken takes an userID and an expiration duration to create a new token
	CreateToken(payType PayloadType, userID int64, duration time.Duration) (string, error)
	// VerifyToken checks if a provided token is valid or not
	VerifyToken(token string) (*Payload, error)
}

/*
FROM AUTH.GO file

// AccessTokenClaims defines the JWT claims for the access token type
type AccessTokenClaims struct {
	UserID  int64
	KeyType string
	jwt.StandardClaims
}

// RefreshTokenClaims defines the JWT claims for the refresh token type
type RefreshTokenClaims struct {
	UserID    int64
	CustomKey string
	KeyType   string
	jwt.StandardClaims
}

*/
