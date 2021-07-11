package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Errors returned by the Create and Validate functions
var (
	ErrExpiredToken = errors.New("Token has expired")
	ErrInvalidToken = errors.New("Invalid token")
)

// Type is a enum to identify the token type
type Type string

const (
	// Access identifies the Access token type
	Access Type = "Access"
	// Refresh identifies the Refresh token type
	Refresh Type = "Refresh"
)

// Payload contains the data of the Token
type Payload struct {
	ID           uuid.UUID `json:"id"`
	TokenType    Type      `json:"tokenType"`
	UserID       int64     `json:"userID"`
	IssuedAt     time.Time `json:"issuedAt"`
	ExpirationAt time.Time `json:"expiredAt"`
}

// NewAccessPayload takes an userId and duration and creates a token payload
func NewAccessPayload(userID int64, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := Payload{
		ID:           tokenID,
		TokenType:    Access,
		UserID:       userID,
		IssuedAt:     time.Now(),
		ExpirationAt: time.Now().Add(duration),
	}
	return &payload, nil
}

// NewRefreshPayload takes an userId and duration and creates a token payload
func NewRefreshPayload(userID int64, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := Payload{
		ID:           tokenID,
		TokenType:    Refresh,
		UserID:       userID,
		IssuedAt:     time.Now(),
		ExpirationAt: time.Now().Add(duration),
	}
	return &payload, nil
}

// Valid checks whether a token is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpirationAt) {
		return ErrExpiredToken
	}
	return nil
}

/* // Custom Key refresh token implementation --
   // Makes the refresh token invalid
// RefreshPayload contains the data of the Refresh Token
type RefreshPayload struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"userID"`
	CustomKey    string    `json:"customKey"`
	IssuedAt     time.Time `json:"issuedAt"`
	ExpirationAt time.Time `json:"expiredAt"`
}

// NOTE: Other option would be to use the jwt.StandardClaims
*/
