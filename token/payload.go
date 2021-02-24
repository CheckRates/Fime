package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Errors returned by the Create and Validate functions
var (
	ErrExpiredToken = errors.New("Token has expired")
	ErrInvalidToken = errors.New("Token is invalid")
)

// PayloadType defines the type of token type
type PayloadType string

const (
	// Refresh specifies the refresh token payload
	Refresh PayloadType = "refresh"
	// Access specifies the access token payload
	Access PayloadType = "access"
)

// Payload contains the data of the token
type Payload struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"userID"`
	TokenType    string    `json:"tokenType"`
	CustomKey    string    `json:"customKey"`
	IssuedAt     time.Time `json:"issuedAt"`
	ExpirationAt time.Time `json:"expiredAt"`
}

// NewPayload takes an userId and duration and creates a token payload
func NewPayload(payType PayloadType, userID int64, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := Payload{
		ID:           tokenID,
		UserID:       userID,
		TokenType:    string(payType),
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
