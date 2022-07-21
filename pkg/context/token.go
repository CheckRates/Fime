package context

import (
	"github.com/checkrates/Fime/pkg/service"
	"github.com/checkrates/Fime/pkg/service/token"
)

// Returns a new instance of the token usecase
func NewTokenUsecase(secret string) (service.TokenMaker, error) {
	maker, err := token.NewJWTMaker(secret)
	if err != nil {
		return nil, err
	}
	return maker, nil
}
