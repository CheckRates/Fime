package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/checkrates/Fime/db/postgres"
	"github.com/labstack/echo"
)

type createUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// createUser takes a JSON request and returns the newly created User object
func (server *Server) createUser(ctx echo.Context) error {
	var req *createUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusOK, errorResponse(err))
	}

	// Make the request to the database and create user
	userArgs := postgres.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := server.store.CreateUser(userArgs)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, user)
}

type getUserParams struct {
	ID int64 `validate:"required,gte=1"`
}

// getUser takes the desired user's ID from the URL and returns a JSON object of requested user
func (server *Server) getUser(ctx echo.Context) error {
	// Parse URL params
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(
			fmt.Errorf("Invalid ID")))
	}

	// Validate the get request params
	req := getUserParams{
		ID: id,
	}

	if err = ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.User(req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, user)
}

type listUserParams struct {
	Page int `validate:"required,min=1"`
	Size int `validate:"required,min=1,max=10"`
}

// listUserParams takes limit and offset params and returns the JSON user objects
func (server *Server) listUsers(ctx echo.Context) error {
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

	// Request list of users to the databse
	arg := postgres.ListParams{
		Limit:  int64(req.Size),
		Offset: int64((req.Page - 1) * req.Size),
	}

	user, err := server.store.Users(arg)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return ctx.JSON(http.StatusOK, user)
}

// TODO: Change to hashpassword
type loginUserParams struct {
	Email    string
	Password string
}

// getUser takes the desired user's ID from the URL and returns a JSON object of requested user
func (server *Server) loginUser(ctx echo.Context) error {
	var req *loginUserParams
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// FIXME: If valid, get user
	if req.Email != "Jon" || req.Password != "secret" {
		return ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("Incorrect Password")))
	}

	u, err := server.store.UserByEmail(req.Email)
	if err != nil {
		// TODO: Better error messages/handling
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// TODO: Check if hashed_password matches stored one in the database
	//if u.HashedPassword != "password" {
	//	return ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("Incorrect password")))
	//}

	// Generates Access and Refresh Tokens

}

func (server *Server) getAccessToken(ctx echo.Context) {
	// FIXME: Get the user from the context
	user := postgres.User{}
	accessToken, err := server.CreateAccessToken(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// FIXME: AccessToken should be a JSON object
	ctx.JSON(http.StatusOK, accessToken)
}
