package service

import (
	"context"
	"time"

	"github.com/checkrates/Fime/pkg/models"
)

type PostUsecase interface {
	Create(ctx context.Context, postData models.PostData) (*models.ImagePost, error)
	FindById(ctx context.Context, id int64) (*models.ImagePost, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, name string, tags []models.CreateTagParams) (*models.ImagePost, error)
	GetMultiple(ctx context.Context, size, page int) ([]models.ImagePost, error)
	GetByUser(ctx context.Context, userId int64, size, page int) ([]models.ImagePost, error)
}

type UserUsecase interface {
	Register(name, email, password string) (*models.UserResponse, error)
	Login(email, password string) (*models.UserResponse, string, error)
	FindById(id int64) (*models.UserResponse, error)
	GetMultiple(size, page int) ([]models.UserResponse, error)
}

type TagUsecase interface {
	GetMultiple(size, page int) ([]models.Tag, error)
	GetUserTags(userId int64, size, page int) ([]models.Tag, error)
}

type TokenMaker interface {
	CreateAccess(userID int64, duration time.Duration) (string, error)
	CreateRefresh(userID int64, duration time.Duration) (string, error)
	VerifyToken(token string) (*models.Payload, error)
}
