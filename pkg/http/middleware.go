package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/checkrates/Fime/pkg/service"
	"github.com/labstack/echo"
)

const (
	authHeaderKey  = "authorization"
	authTypeBearer = "bearer"
	authPayloadKey = "auth_payload"
)

func AuthMiddleware(tokenMaker service.TokenMaker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authHeader := ctx.Request().Header.Get(authHeaderKey)
			if len(authHeader) == 0 {
				err := errors.New("authorization header is not provided")
				return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			}

			fields := strings.Fields(authHeader)
			if len(fields) < 2 {
				err := errors.New("invalid authorization header format")
				return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			}

			// Handle auth based on type retrieved from header (May implement OAuth later)
			authType := strings.ToLower(fields[0])
			if authType != authTypeBearer {
				err := errors.New("provided authorization type not supported by the server")
				return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			}

			accessToken := fields[1]
			payload, err := tokenMaker.VerifyToken(accessToken)
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			}

			ctx.Set(authPayloadKey, payload)

			return next(ctx)
		}
	}
}

// errorResponse formats error to Echo
func errorResponse(err error) echo.Map {
	return echo.Map{"error": err.Error()}
}
