package user

import (
	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/checkrates/Fime/pkg/storage"
	"github.com/checkrates/Fime/util"
)

type userService struct {
	repo   storage.UserRepository
	token  service.TokenMaker
	config config.Config
}

func NewUserService(user storage.UserRepository, token service.TokenMaker, config config.Config) service.UserUsecase {
	return userService{
		repo:   user,
		token:  token,
		config: config,
	}
}

// Register a new user into Fime and returns the newly created User object
func (u userService) Register(name, email, password string) (*models.UserResponse, error) {
	// Hash password before saving to the repository
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	userArgs := models.CreateUserParams{
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
	}

	newUser, err := u.repo.Create(userArgs)
	if err != nil {
		return nil, err
	}

	resp := models.NewUserResponse(newUser)
	return &resp, nil
}

// Login a user with email and password. Returns the user object and access token if successful
func (u userService) Login(email, password string) (*models.UserResponse, string, error) {
	retUser, err := u.repo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}

	err = util.ValidatePassword(password, retUser.HashedPassword)
	if err != nil {
		return nil, "", err
	}

	accessToken, err := u.token.CreateAccess(retUser.ID, u.config.Token.AccessExpiration)
	if err != nil {
		return nil, "", err
	}

	resp := models.NewUserResponse(retUser)
	return &resp, accessToken, nil
}

// Returns an user by id, if found
func (u userService) FindById(id int64) (*models.UserResponse, error) {
	user, err := u.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	resp := models.NewUserResponse(user)
	return &resp, nil
}

// Takes a size and the page number to provide a subset of users
func (u userService) GetMultiple(size, page int) ([]models.UserResponse, error) {
	arg := models.ListUserParams{
		Limit:  size,
		Offset: (page - 1) * size,
	}

	users, err := u.repo.GetMultiple(arg)
	if err != nil {
		return nil, err
	}

	var resp []models.UserResponse
	for _, u := range users {
		resp = append(resp, models.NewUserResponse(u))
	}

	return resp, nil
}
