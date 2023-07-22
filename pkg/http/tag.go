package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/labstack/echo"
)

type TagPort interface {
	GetMultiple(ctx echo.Context) error
	GetUserTags(ctx echo.Context) error
}

type tagApi struct {
	tag service.TagUsecase
}

// Returns the default implementation of the tag port.
func NewTagApi(tag service.TagUsecase) TagPort {
	return tagApi{
		tag: tag,
	}
}

type listTagParams struct {
	Page int `validate:"required,min=1"`
	Size int `validate:"required,min=1,max=10"`
}

// Takes limit and offset params and returns the JSON tags objects
func (t tagApi) GetMultiple(ctx echo.Context) error {
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

	// Validates the URL params, they are grouped in a struct to facilitate the use
	// of the Echo#Valicator
	req := listTagParams{
		Page: page,
		Size: size,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	authPayload := ctx.Get(authPayloadKey).(*models.Payload)
	if authPayload == nil {
		return ctx.JSON(http.StatusUnauthorized, errorResponse(
			fmt.Errorf("cannot access tags without being logged in")))
	}

	tags, err := t.tag.GetMultiple(req.Size, req.Page)
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

// Takes limit and offset params and returns all user tags
func (t tagApi) GetUserTags(ctx echo.Context) error {
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

	// Validates the URL params, they are grouped in a struct to facilitate the use
	// of the Echo#Valicator
	req := listUserTagsParams{
		ID:   id,
		Page: page,
		Size: size,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	authPayload := ctx.Get(authPayloadKey).(*models.Payload)
	if req.ID != authPayload.UserID {
		return ctx.JSON(http.StatusUnauthorized, errorResponse(
			fmt.Errorf("cannot access user tags without being logged in")))
	}

	tags, err := t.tag.GetUserTags(req.ID, req.Size, req.Page)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, tags)
}
