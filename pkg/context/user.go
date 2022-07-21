package context

import (
	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/pkg/http"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/checkrates/Fime/pkg/service/token"
	"github.com/checkrates/Fime/pkg/service/user"
	"github.com/checkrates/Fime/pkg/storage/db/postgres"
	"github.com/jmoiron/sqlx"
)

// Returns a new instance of the user usecase
func NewUserUsecase(db *sqlx.DB, config config.Config) (service.UserUsecase, error) {
	tokenMaker, err := token.NewJWTMaker(config.Token.AccessSecret)
	if err != nil {
		return nil, err
	}

	return user.NewUserService(
		postgres.NewUserRepository(db),
		tokenMaker,
		config,
	), nil
}

// Returns a new HTTP implementation of the userPort
func NewUserPort(db *sqlx.DB, config config.Config) (http.UserPort, error) {
	usecase, err := NewUserUsecase(db, config)
	if err != nil {
		return nil, err
	}
	return http.NewUserApi(usecase), nil
}
