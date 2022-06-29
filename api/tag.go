package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/checkrates/Fime/db/postgres"
	"github.com/labstack/echo"
)

// listTags takes limit and offset params and returns the JSON tags objects
func (server *Server) listTags(ctx echo.Context) error {
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("invalid page value")))
	}

	size, err := strconv.Atoi(ctx.QueryParam("size"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("invalid size value")))
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

type listUserTagsParams struct {
	ID   int64 `validate:"required,min=1"`
	Page int   `validate:"required,min=1"`
	Size int   `validate:"required,min=1,max=10"`
}

// listTags takes limit and offset params and returns all user tags
func (server *Server) listUserTags(ctx echo.Context) error {
	// Parse URL params
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("invalid ID")))
	}

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("invalid page value")))
	}

	size, err := strconv.Atoi(ctx.QueryParam("size"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("invalid size value")))
	}

	// Validate list request params
	req := listUserTagsParams{
		ID:   id,
		Page: page,
		Size: size,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Request list of tags to the databse
	arg := postgres.ListUserTagsParams{
		ID:     req.ID,
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	}

	tags, err := server.store.GetUserTags(arg)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, tags)
}
