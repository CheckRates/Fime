package post

import (
	"context"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/checkrates/Fime/pkg/storage"
)

type postService struct {
	repo   storage.PostRepository
	bucket service.BucketUsecase
}

func NewPostService(post storage.PostRepository, bucket service.BucketUsecase) service.PostUsecase {
	return postService{
		repo:   post,
		bucket: bucket,
	}
}

func (p postService) Create(ctx context.Context, postData models.PostData) (*models.ImagePost, error) {
	imgMeta, err := p.bucket.RequestUpload(models.RequestUploadParams{
		Filename:   postData.Name,
		UserID:     postData.UserId,
		Fileheader: postData.EncodedImg,
	})
	if err != nil {
		return nil, err
	}

	arg := models.CreatePostParams{
		Name:   postData.Name,
		URL:    imgMeta.ImageKey,
		UserID: postData.UserId,
		Tags:   postData.Tags,
	}

	imgPost, err := p.repo.Create(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &imgPost, nil
}

func (p postService) FindById(ctx context.Context, id int64) (*models.ImagePost, error) {
	imgPost, err := p.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &imgPost, nil
}

func (p postService) Delete(ctx context.Context, id int64) error {
	imgPost, err := p.repo.FindById(ctx, id)
	if err != nil {
		return err
	}

	err = p.bucket.Delete(imgPost.Image.URL)
	if err != nil {
		return err
	}

	err = p.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (p postService) Update(ctx context.Context, id int64, name string, tags []models.CreateTagParams) (*models.ImagePost, error) {
	arg := models.UpdatePostParams{
		ID:   id,
		Name: name,
		Tags: tags,
	}

	imgPost, err := p.repo.Update(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &imgPost, nil
}

func (p postService) GetMultiple(ctx context.Context, size, page int) ([]models.ImagePost, error) {
	arg := models.ListImagesParams{
		Limit:  size,
		Offset: (page - 1) * size,
	}

	imgs, err := p.repo.GetMutiple(ctx, arg)
	if err != nil {
		return nil, err
	}

	return imgs, nil
}

func (p postService) GetByUser(ctx context.Context, userId int64, size, page int) ([]models.ImagePost, error) {
	arg := models.ListUserImagesParams{
		UserID: userId,
		Limit:  size,
		Offset: (page - 1) * size,
	}

	imgs, err := p.repo.GetByUser(ctx, arg)
	if err != nil {
		return nil, err

	}
	return imgs, nil
}
