package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/checkrates/Fime/db/postgres"
	"github.com/labstack/echo"
)

type createTagParams struct {
	Tag string
}

type postImageParams struct {
	Name       string                     `json:"name"`
	EncodedImg string                     `json:"image"`
	UserID     int64                      `json:"ownerID"`
	Tags       []postgres.CreateTagParams `json:"tags"`
}

func (server *Server) postImage(ctx echo.Context) error {
	var req *postImageParams
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Upload image to S3 bucket and get resource URL
	imgURL, err := server.UploadImage(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse((err)))
	}

	// Make the request to the database and post image
	arg := postgres.MakePostParams{
		Name:   req.Name,
		URL:    imgURL,
		UserID: req.UserID,
		Tags:   req.Tags,
	}

	imgs, err := server.store.MakePostTx(ctx.Request().Context(), arg)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, imgs)
}

type getImageParam struct {
	ID int64 `validate:"required,min=1"`
}

func (server *Server) getImage(ctx echo.Context) error {
	// Parse URL params
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("Invalid ID")))
	}

	// Validate the get request params
	req := getImageParam{
		ID: id,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	img, err := server.store.GetPostTx(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, img)
}

type deleteImageParam struct {
	ID int64 `validate:"required,min=1"`
}

func (server *Server) deleteImage(ctx echo.Context) error {
	// Parse URL params
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("Invalid ID")))
	}

	// Validate the delete request params
	req := deleteImageParam{
		ID: id,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Delete image in the S3 repo
	if err = server.DeleteImage(req.ID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	err = server.store.DeletePostTx(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, "Image ID "+fmt.Sprint(req.ID)+" deleted")
}

type updatePostParams struct {
	ID   int64
	Name string
	Tags []postgres.CreateTagParams
}

func (server *Server) updateImage(ctx echo.Context) error {
	var req *updatePostParams
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusOK, errorResponse(err))
	}

	// Make the request to the database and update image post
	arg := postgres.UpdatePostParams{
		ID:   req.ID,
		Name: req.Name,
		Tags: req.Tags,
	}

	imgPost, err := server.store.UpdatePostTx(ctx.Request().Context(), arg)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, imgPost)
}

type listPostParams struct {
	Page int `validate:"required,min=1"`
	Size int `validate:"required,min=1,max=10"`
}

func (server *Server) listImage(ctx echo.Context) error {
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
	req := listPostParams{
		Page: page,
		Size: size,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Request list of image posts to the databse
	arg := postgres.ListParams{
		Limit:  int64(req.Size),
		Offset: int64((req.Page - 1) * req.Size),
	}

	imgs, err := server.store.ListPostTx(ctx.Request().Context(), arg)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, imgs)
}

type listUserPostsParams struct {
	UserID int64 `validate:"required,min=1"`
	Page   int   `validate:"required,min=1"`
	Size   int   `validate:"required,min=1,max=10"`
}

func (server *Server) listUserImages(ctx echo.Context) error {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("Invalid ID value")))
	}

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
	req := listUserPostsParams{
		UserID: userID,
		Page:   page,
		Size:   size,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Request list of user's image posts to the database
	arg := postgres.ListUserPostsParams{
		UserID: userID,
		Limit:  int64(req.Size),
		Offset: int64((req.Page - 1) * req.Size),
	}

	imgs, err := server.store.ListUserPostTx(ctx.Request().Context(), arg)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, imgs)
}
