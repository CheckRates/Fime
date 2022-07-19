package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/labstack/echo"
)

// errorResponse formats error to Echo
func errorResponse(err error) echo.Map {
	return echo.Map{"error": err.Error()}
}

type PostPort interface {
	Create(ctx echo.Context) error
	FindById(ctx echo.Context) error
	Delete(ctx echo.Context) error
	Update(ctx echo.Context) error
	GetMultiple(ctx echo.Context) error
	GetByUserId(ctx echo.Context) error
}

type postApi struct {
	post service.PostUsecase
}

// Returns the default implementation of the tag port.
func NewPostApi(post service.PostUsecase) PostPort {
	return postApi{
		post: post,
	}
}

type postImageRequest struct {
	Name       string                   `json:"name"`
	EncodedImg string                   `json:"image"`
	Tags       []models.CreateTagParams `json:"tags"`
}

func (p postApi) Create(ctx echo.Context) error {
	var req *postImageRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Retrieve the authenticated user
	authPayload := ctx.Get(authPayloadKey).(*models.Payload)
	if authPayload == nil {
		err := errors.New("cannot retrived user if not authenticated")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	imgPost, err := p.post.Create(ctx.Request().Context(), models.PostData{
		Name:       req.Name,
		EncodedImg: req.EncodedImg,
		UserId:     authPayload.UserID,
		Tags:       req.Tags,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, imgPost)
}

type getImageParam struct {
	ID int64 `validate:"required,min=1"`
}

func (p postApi) FindById(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("invalid ID")))
	}

	req := getImageParam{
		ID: id,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	imgPost, err := p.post.FindById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// If authenticated user is not the owner of the post, return error
	authPayload := ctx.Get(authPayloadKey).(*models.Payload)
	if imgPost.Image.OwnerID != authPayload.UserID {
		err := errors.New("image does not belong to authenticated user")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, imgPost)
}

type deleteImageParam struct {
	ID int64 `validate:"required,min=1"`
}

func (p postApi) Delete(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("invalid ID")))
	}

	req := deleteImageParam{
		ID: id,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	imgPost, err := p.post.FindById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	authPayload := ctx.Get(authPayloadKey).(*models.Payload)
	if imgPost.Image.OwnerID != authPayload.UserID {
		err := errors.New("image does not belong to authenticated user")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	err = p.post.Delete(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, "Image ID "+fmt.Sprint(req.ID)+" deleted")
}

type updatePostRequest struct {
	ID   int64
	Name string
	Tags []models.CreateTagParams
}

func (p postApi) Update(ctx echo.Context) error {
	var req *updatePostRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusOK, errorResponse(err))
	}

	// Check if correct user is authenticated
	imgPost, err := p.post.FindById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	authPayload := ctx.Get(authPayloadKey).(*models.Payload)
	if imgPost.Image.OwnerID != authPayload.UserID {
		err := errors.New("image does not belong to authenticated user")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	imgPost, err = p.post.Update(ctx.Request().Context(), req.ID, req.Name, req.Tags)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, imgPost)
}

type listPostParams struct {
	Page int `validate:"required,min=1"`
	Size int `validate:"required,min=1,max=10"`
}

func (p postApi) GetMultiple(ctx echo.Context) error {
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

	req := listPostParams{
		Page: page,
		Size: size,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	authPayload := ctx.Get(authPayloadKey).(*models.Payload)
	if authPayload == nil {
		return ctx.JSON(http.StatusUnauthorized, errorResponse(
			fmt.Errorf("cannot access images without being logged in")))
	}

	imgPosts, err := p.post.GetMultiple(ctx.Request().Context(), req.Size, req.Page)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, imgPosts)
}

type listUserPostsParams struct {
	UserID int64 `validate:"required,min=1"`
	Page   int   `validate:"required,min=1"`
	Size   int   `validate:"required,min=1,max=10"`
}

func (p postApi) GetByUserId(ctx echo.Context) error {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("invalid ID value")))
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

	req := listUserPostsParams{
		UserID: userID,
		Page:   page,
		Size:   size,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	authPayload := ctx.Get(authPayloadKey).(*models.Payload)
	if req.UserID != authPayload.UserID {
		err := errors.New("images do not belong to authenticated user")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	imgs, err := p.post.GetByUser(ctx.Request().Context(), req.UserID, req.Size, req.Page)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, imgs)
}
