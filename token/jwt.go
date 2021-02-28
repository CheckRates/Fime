package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretSize = 32

// JWTMaker is a a JWT constructor
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretSize {
		return nil, fmt.Errorf("Invalid secret key, size must be at least %d digits long", minSecretSize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken takes an userID and an expiration duration to create a new token
func (maker *JWTMaker) CreateToken(payType PayloadType, userID int64, duration time.Duration) (string, error) {

	payload, err := NewPayload(payType, userID, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	return token.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if a provided token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	// keyFunc needs to return the signing key of the JWT token
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidToken
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(maker.secretKey))
		if err != nil {
			return nil, err
		}
		return verifyKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		// Check which type of error was returned from the ParseWithClaims func
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Finally, cast into a payload object and returned it
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}