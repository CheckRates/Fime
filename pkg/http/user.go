package http

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserPort interface {
	Register(ctx echo.Context) error
	Login(ctx echo.Context) error
	FindById(ctx echo.Context) error
	GetMultiple(ctx echo.Context) error
}

type userApi struct {
	user service.UserUsecase
}

// Returns the default implementation of the user port.
func NewUserApi(user service.UserUsecase) UserPort {
	return userApi{
		user: user,
	}
}

type registerUserRequest struct {
	Name     string `json:"name" validate:"required,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Takes a JSON request and returns the newly created User object
func (u userApi) Register(ctx echo.Context) error {
	var req *registerUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusOK, errorResponse(err))
	}

	user, err := u.user.Register(req.Name, req.Email, req.Password)
	if err != nil {
		// Possible Database errors
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return ctx.JSON(http.StatusForbidden, errorResponse(err))
			}
		}

		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusCreated, user)
}

type loginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginUserResponse struct {
	AccessToken string              `json:"access_token"`
	User        models.UserResponse `json:"user"`
}

func (u userApi) Login(ctx echo.Context) error {
	var req *loginUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, accessToken, err := u.user.Login(req.Email, req.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else if err == bcrypt.ErrMismatchedHashAndPassword {
			return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		}
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	resp := loginUserResponse{
		AccessToken: accessToken,
		User:        *user,
	}
	return ctx.JSON(http.StatusOK, resp)
}

type getUserParams struct {
	ID int64 `validate:"required,gte=1"`
}

// Takes the desired user's ID from the URL and returns a JSON object of requested user
func (u userApi) FindById(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("invalid ID")))
	}

	req := getUserParams{
		ID: id,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	authPayload := ctx.Get(authPayloadKey).(*models.Payload)
	if authPayload == nil {
		err := errors.New("cannot retrived user if not authenticated")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	user, err := u.user.FindById(req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, user)
}

type listUserParams struct {
	Page int `validate:"required,min=1"`
	Size int `validate:"required,min=1,max=10"`
}

// Takes limit and offset params and returns the JSON user objects
func (u userApi) GetMultiple(ctx echo.Context) error {
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

	req := listUserParams{
		Page: page,
		Size: size,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	authPayload := ctx.Get(authPayloadKey).(*models.Payload)
	if authPayload == nil {
		err := errors.New("cannot retrived users if not authenticated")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	users, err := u.user.GetMultiple(req.Size, req.Page)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, users)
}
