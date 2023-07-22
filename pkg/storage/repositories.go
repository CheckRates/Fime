package storage

import (
	"context"

	"github.com/checkrates/Fime/pkg/models"
)

type PostRepository interface {
	FindById(ctx context.Context, id int64) (models.ImagePost, error)
	GetMutiple(ctx context.Context, arg models.ListImagesParams) ([]models.ImagePost, error)
	GetByUser(ctx context.Context, arg models.ListUserImagesParams) ([]models.ImagePost, error)
	Create(ctx context.Context, arg models.CreatePostParams) (models.ImagePost, error)
	Update(ctx context.Context, arg models.UpdatePostParams) (models.ImagePost, error)
	Delete(ctx context.Context, id int64) error
}

type UserRepository interface {
	FindById(id int64) (models.User, error)
	FindByEmail(email string) (models.User, error)
	GetMultiple(args models.ListUserParams) ([]models.User, error)
	Create(args models.CreateUserParams) (models.User, error)
	Update(args models.UpdateUserParams) (models.User, error)
	Delete(id int64) error
}

type ImageRepository interface {
	FindById(id int64) (models.Image, error)
	GetMultiple(args models.ListImagesParams) ([]models.Image, error)
	GetByUser(args models.ListUserImagesParams) ([]models.Image, error)
	Create(args models.CreateImageParams) (models.Image, error)
	Update(args models.UpdateImageParams) (models.Image, error)
	Delete(id int64) error
}

type TagRepository interface {
	FindById(id int64) (models.Tag, error)
	GetMultiple(args models.ListTagsParams) ([]models.Tag, error)
	Create(args models.CreateTagParams) (models.Tag, error)
	GetUserTags(arg models.ListUserTagsParams) ([]models.Tag, error)
	Delete(id int64) error
}

type ImageTagRepository interface {
	Create(it models.ImageTag) error
	Delete(it models.ImageTag) error
	DeleteAllFromImage(imgID int64) error
	GetTagsByImageID(imgID int64) ([]models.Tag, error)
	GetImagesByTagID(tagID int64) ([]models.Image, error)
}
