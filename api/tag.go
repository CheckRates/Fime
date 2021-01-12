package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/checkrates/Fime/db/postgres"
	"github.com/labstack/echo"
)

type lisTagParams struct {
	Page int `validate:"required,min=1"`
	Size int `validate:"required,min=1,max=10"`
}

// listTags takes limit and offset params and returns the JSON tags objects
func (server *Server) listTags(ctx echo.Context) error {
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("Invalid page value")))
	}

	size, err := strconv.Atoi(ctx.QueryParam("size"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("Invalid size value")))
	}

	// Validate list request params
	req := listUserParams{
		Page: page,
		Size: size,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Request list of tags to the databse
	arg := postgres.ListParams{
		Limit:  int64(req.Size),
		Offset: int64((req.Page - 1) * req.Size),
	}

	tags, err := server.store.Tags(arg)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, tags)
}
