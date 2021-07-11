package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// MiddlewareValidateRefreshToken takes a refresh token and checks if it is a valid token
func (server *Server) MiddlewareValidateRefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, err := extractToken(ctx.Request().Header)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("error")))
			return nil
		}

		_, err = server.token.VerifyToken(token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return nil
		}

		/*
			//Custom Key Validation

			user, err := server.store.User(payload.UserID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return nil
			}

			actualCustomKey := uh.authService.GenerateCustomKey(user.ID, user.TokenHash)
			if customKey != actualCustomKey {
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return nil
			}

			ctx := context.WithValue(r.Context(), UserKey{}, *user)
			r = r.WithContext(ctx)

			ctx.Set("user", user)
		*/
		return next(ctx)
	}
}

// MiddlewareValidateAccessToken checks if a access token and checks if it is valid token
func (server *Server) MiddlewareValidateAccessToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, err := extractToken(ctx.Request().Header)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("error")))
			return nil
		}

		payload, err := server.token.VerifyToken(token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return nil
		}

		ctx.Set("userID", payload.UserID)
		return next(ctx)
	}
}

// extractToken returns the token from the HTTP header if it exist
func extractToken(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	authHeaderContent := strings.Split(authHeader, " ")
	if len(authHeaderContent) != 2 {
		return "", errors.New("Token not provided or malformed")
	}
	return authHeaderContent[1], nil
}
