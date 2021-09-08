package api

import (
	"net/http"

	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/util"
	"github.com/labstack/echo"
)

type loginUserResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type loginUserParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// loginUser takes an user email and password and returns a access and refresh token,
// if the user is valid
func (server *Server) loginUser(ctx echo.Context) error {
	var req *loginUserParams
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Get user from db and check if credentials match
	user, err := server.store.UserByEmail(req.Email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	err = util.ValidatePassword(req.Password, user.HashedPassword)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	// User successfully login -- Generates Access and Refresh Tokens
	accessToken, err := server.token.CreateAccess(
		user.ID,
		config.New().Token.AccessExpiration,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	refreshToken, err := server.token.CreateRefresh(
		user.ID,
		config.New().Token.RefreshExpiration,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, &loginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
