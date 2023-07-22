package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("invalid token")
)

// Payload contains the data of the token
type Payload struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"userID"`
	IssuedAt     time.Time `json:"issuedAt"`
	ExpirationAt time.Time `json:"expiredAt"`
}

// Takes an userId and duration and creates a token payload
func NewAccessPayload(userID int64, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := Payload{
		ID:           tokenID,
		UserID:       userID,
		IssuedAt:     time.Now(),
		ExpirationAt: time.Now().Add(duration),
	}
	return &payload, nil
}

// Takes an userId and duration and creates a token payload
func NewRefreshPayload(userID int64, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := Payload{
		ID:           tokenID,
		UserID:       userID,
		IssuedAt:     time.Now(),
		ExpirationAt: time.Now().Add(duration),
	}
	return &payload, nil
}

// Valid checks whether a token is expired or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpirationAt) {
		return ErrExpiredToken
	}
	return nil
}
