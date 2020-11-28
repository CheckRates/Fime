package api

import (
	"net/http"

	"github.com/checkrates/Fime/fime"
	"github.com/labstack/echo"
)

type createUserRequest struct {
	Name string `json:"name" validate:"required"`
}

func (server *Server) createUser(ctx echo.Context) error {
	/*
		var req createUserRequest
		if err := ctx.Bind(req); err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse(err))
		}
	*/

	user := fime.User{
		Name: "Dab",
	}

	err := server.store.CreateUser(&user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, user)
}
