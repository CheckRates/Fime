package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/db/postgres"
	"github.com/checkrates/Fime/util"
	"github.com/labstack/echo"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Name     string `json:"name" validate:"required,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type createUserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
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

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// Make the request to the database and create user
	userArgs := postgres.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(userArgs)
	if err != nil {
		// Possible Database errors
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return ctx.JSON(http.StatusForbidden, errorResponse(err))
			}
		}
	}

	resp := createUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	return ctx.JSON(http.StatusOK, resp)
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

type loginUserParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginUserResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

	// User successfully logd Generates Access and Refresh Tokens
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

/*
func (server *Server) getAccessToken(ctx echo.Context) {
	// FIXME: Get the user from the context
	user := postgres.User{}
	accessToken, err := server.token.CreateToken(token.Access, user.ID, config.New().Token.AccessExpiration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	// FIXME: AccessToken should be a JSON object
	ctx.JSON(http.StatusOK, accessToken)
}
*/
