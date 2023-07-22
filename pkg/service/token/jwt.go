package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/dgrijalva/jwt-go"
)

// TODO: Change to 32 later
const minSecretSize = 3

type jwtMaker struct {
	secretKey string
}

// Creates a new JWTMaker
func NewJWTMaker(secretKey string) (service.TokenMaker, error) {
	if len(secretKey) < minSecretSize {
		return nil, fmt.Errorf("invalid secret key, size must be at least %d digits long", minSecretSize)
	}
	return &jwtMaker{secretKey}, nil
}

// Takes an userID and an expiration duration to create a new access token. This
// short lived token will be used in every request that requires user authentication
func (maker *jwtMaker) CreateAccess(userID int64, duration time.Duration) (string, error) {
	payload, err := models.NewAccessPayload(userID, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(maker.secretKey))
}

// Takes an userID and expiration to create a longer lived token. This token is used
// to request a new access token for a valid user
func (maker *jwtMaker) CreateRefresh(userID int64, duration time.Duration) (string, error) {
	payload, err := models.NewRefreshPayload(userID, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	return token.SignedString([]byte(maker.secretKey))
}

// Checks if a provided token is valid or not
func (maker *jwtMaker) VerifyToken(token string) (*models.Payload, error) {
	// keyFunc checks if the received token has the expected signing method
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &models.Payload{}, keyFunc)
	if err != nil {
		// Check which type of error was returned from the ParseWithClaims func
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, models.ErrExpiredToken) {
			return nil, models.ErrExpiredToken
		}
		return nil, models.ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*models.Payload)
	if !ok {
		return nil, models.ErrInvalidToken
	}

	return payload, nil
}
