package post

import (
	"context"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/checkrates/Fime/pkg/storage"
)

type PostService struct {
	post storage.PostRepository
}

func NewPostService(post storage.PostRepository) service.PostUsecase {
	return PostService{
		post: post,
	}
}

func (p PostService) Create(ctx context.Context, postData models.PostData) (*models.ImagePost, error) {
	// Upload image to S3 bucket and get resource URL
	//imgURL, err := server.UploadImage(encondedImg)
	//if err != nil {
	//	return ctx.JSON(http.StatusInternalServerError, errorResponse((err)))
	//}
	imgURL := "www.coolimage.com" // FIXME: Connect to AWS S3 bucket

	arg := models.CreatePostParams{
		Name:   postData.Name,
		URL:    imgURL,
		UserID: postData.UserId,
		Tags:   postData.Tags,
	}

	imgPost, err := p.post.Create(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &imgPost, nil
}

func (p PostService) FindById(ctx context.Context, id int64) (*models.ImagePost, error) {
	imgPost, err := p.post.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &imgPost, nil
}

func (p PostService) Delete(ctx context.Context, id int64) error {
	_, err := p.post.FindById(ctx, id)
	if err != nil {
		return err
	}

	// FIXME:
	// Delete image in the S3 repo
	// if err = server.DeleteImage(req.ID); err != nil {
	//	return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//}

	err = p.post.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (p PostService) Update(ctx context.Context, id int64, name string, tags []models.CreateTagParams) (*models.ImagePost, error) {
	arg := models.UpdatePostParams{
		ID:   id,
		Name: name,
		Tags: tags,
	}

	imgPost, err := p.post.Update(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &imgPost, nil
}

func (p PostService) GetMultiple(ctx context.Context, size, page int) ([]models.ImagePost, error) {
	arg := models.ListImagesParams{
		Limit:  size,
		Offset: (page - 1) * size,
	}

	imgs, err := p.post.GetMutiple(ctx, arg)
	if err != nil {
		return nil, err
	}

	return imgs, nil
}

func (p PostService) GetByUser(ctx context.Context, userId int64, size, page int) ([]models.ImagePost, error) {
	arg := models.ListUserImagesParams{
		UserID: userId,
		Limit:  size,
		Offset: (page - 1) * size,
	}

	imgs, err := p.post.GetByUser(ctx, arg)
	if err != nil {
		return nil, err

	}
	return imgs, nil
}
