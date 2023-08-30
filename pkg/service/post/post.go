package post

import (
	"context"
	"fmt"
	"math"

	"github.com/checkrates/Fime/pkg/filetype"
	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/checkrates/Fime/pkg/storage"
)

type postService struct {
	repo   storage.PostRepository
	bucket service.BucketUsecase
}

// in Megabytes
const sizeImageChunk = 5000 // TODO: Change

func NewPostService(post storage.PostRepository, bucket service.BucketUsecase) service.PostUsecase {
	return postService{
		repo:   post,
		bucket: bucket,
	}
}

func (p postService) RequestUpload(ctx context.Context, req models.RequestUploadParams) (*models.RequestUploadResponse, error) {
	valid, err := filetype.IsValid(req.EncodedImgHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file")
	}
	if !valid {
		return nil, fmt.Errorf("image file type not supported")
	}

	imgMeta, err := p.bucket.InitiateUpload(ctx, models.InitiateUploadParams{
		Filename:   req.Filename,
		UserID:     req.UserId,
		Fileheader: req.EncodedImgHeader,
	})
	if err != nil {
		return nil, err
	}

	// Determine number of image chunks, therefore number of presigned URLs necessary
	// to upload the requested image
	numURLs := math.Ceil(float64(req.ImageSize) / sizeImageChunk)

	return &models.RequestUploadResponse{
		Filename: req.Filename,
		NumParts: int(numURLs),
		UploadId: imgMeta.UploadId,
		ImageKey: imgMeta.ImageKey,
	}, nil
}

func (p postService) GetUploadURLs(ctx context.Context, uploadId string, numParts int32) ([]models.PresignedURL, error) {
	urls, err := p.bucket.GeneratePresignURLs(ctx, uploadId, numParts)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (p postService) CompleteUpload(ctx context.Context, req models.CompleteUploadParams) (*models.ImagePost, error) {
	url, err := p.bucket.CompleteUpload(ctx, req.UploadID)
	if err != nil {
		return nil, err
	}

	arg := models.CreatePostParams{
		Filename: req.Filename,
		URL:      url,
		UserID:   req.UserID,
		Tags:     req.Tags,
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

	err = p.bucket.Delete(ctx, imgPost.Image.URL)
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
