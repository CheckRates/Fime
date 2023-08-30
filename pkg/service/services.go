package service

import (
	"context"
	"time"

	"github.com/checkrates/Fime/pkg/models"
)

type PostUsecase interface {
	RequestUpload(ctx context.Context, req models.RequestUploadParams) (*models.RequestUploadResponse, error)
	GetUploadURLs(ctx context.Context, uploadId string, numParts int32) ([]models.PresignedURL, error)
	CompleteUpload(ctx context.Context, req models.CompleteUploadParams) (*models.ImagePost, error)
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

type BucketUsecase interface {
	InitiateUpload(ctx context.Context, uploadParams models.InitiateUploadParams) (*models.InitiateUploadResponse, error)
	GeneratePresignURLs(ctx context.Context, uploadId string, numParts int32) ([]models.PresignedURL, error)
	CompleteUpload(ctx context.Context, uploadId string) (string, error)
	AbortUpload(ctx context.Context, imageKey string, uploadId string) error
	Delete(ctx context.Context, key string) error
}

type TokenMaker interface {
	CreateAccess(userID int64, duration time.Duration) (string, error)
	CreateRefresh(userID int64, duration time.Duration) (string, error)
	VerifyToken(token string) (*models.Payload, error)
}
