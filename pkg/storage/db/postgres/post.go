package postgres

import (
	"context"
	"errors"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/jmoiron/sqlx"
)

type PostSQL struct {
	db       *sqlx.DB
	image    *ImageSQL
	tag      *TagSQL
	imageTag *ImageTagSQL
}

// Returns an access point to Fime's image Posts (image + tags)
func NewPostRepository(db *sqlx.DB) *PostSQL {
	return &PostSQL{
		db:       db,
		image:    NewImageRepository(db),
		tag:      NewTagRepository(db),
		imageTag: NewImageTagRepository(db),
	}
}

// Helper function used for execting transactions
func (s *PostSQL) execTx(ctx context.Context, fn func(*PostSQL) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(s)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			allErrors := "Transaction Error: " + rbErr.Error() + "\n" + err.Error()
			return errors.New(allErrors)
		}
		return err
	}

	return tx.Commit()
}

// Creates an image post in Fime. The post contains related tags and the information about the image (title, URL)
func (s *PostSQL) Create(ctx context.Context, arg models.CreatePostParams) (models.ImagePost, error) {
	var retPost models.ImagePost

	err := s.execTx(ctx, func(s *PostSQL) error {
		var err error

		imageArgs := models.CreateImageParams{
			Name:    arg.Name,
			URL:     arg.URL,
			OwnerID: arg.UserID,
		}
		img, err := s.image.Create(imageArgs)
		if err != nil {
			return err
		}

		tags, err := s.tagImagePost(img.ID, arg.Tags)
		if err != nil {
			return err
		}

		retPost.Image = img
		retPost.Tags = tags
		return nil
	})

	return retPost, err
}

// Gets a image post from the database, if found
func (s *PostSQL) FindById(ctx context.Context, id int64) (models.ImagePost, error) {
	var retPost models.ImagePost

	err := s.execTx(ctx, func(s *PostSQL) error {
		var err error

		img, err := s.image.FindById(id)
		if err != nil {
			return err
		}

		tags, err := s.imageTag.GetTagsByImageID(img.ID)
		if err != nil {
			return err
		}

		// No errors -> proceed to return post
		retPost.Image = img
		retPost.Tags = tags
		return nil
	})

	// Return created post and error
	return retPost, err
}

// Gets a subset of image posts from the database
func (s *PostSQL) GetMutiple(ctx context.Context, arg models.ListImagesParams) ([]models.ImagePost, error) {
	var retPost []models.ImagePost

	err := s.execTx(ctx, func(s *PostSQL) error {
		var err error

		imgs, err := s.image.GetMultiple(arg)
		if err != nil {
			return err
		}

		// Get every image posts tags
		for _, img := range imgs {
			tags, err := s.imageTag.GetTagsByImageID(img.ID)
			if err != nil {
				return err
			}

			post := models.ImagePost{
				Image: img,
				Tags:  tags,
			}

			retPost = append(retPost, post)
		}

		return nil
	})

	// Return the list of posts
	return retPost, err
}

// Gets a list of all the user's image posts from the database
func (s *PostSQL) GetByUser(ctx context.Context, arg models.ListUserImagesParams) ([]models.ImagePost, error) {
	var retPost []models.ImagePost

	err := s.execTx(ctx, func(s *PostSQL) error {
		var err error

		imgs, err := s.image.GetByUser(arg)
		if err != nil {
			return err
		}

		// No images post case
		if len(imgs) == 0 {
			return nil
		}

		// Get every image posts tags
		for _, img := range imgs {
			tags, err := s.imageTag.GetTagsByImageID(img.ID)
			if err != nil {
				return err
			}

			post := models.ImagePost{
				Image: img,
				Tags:  tags,
			}

			retPost = append(retPost, post)
		}

		return nil
	})

	// Return the list of posts
	return retPost, err
}

// Updates image post's tags info and image name
func (s *PostSQL) Update(ctx context.Context, arg models.UpdatePostParams) (models.ImagePost, error) {
	var retPost models.ImagePost

	err := s.execTx(ctx, func(s *PostSQL) error {
		var err error

		// Remove all tags from post and (re)add updated tags
		err = s.imageTag.DeleteAllFromImage(arg.ID)
		if err != nil {
			return err
		}
		tags, err := s.tagImagePost(arg.ID, arg.Tags)
		if err != nil {
			return err
		}

		updateImg := models.UpdateImageParams{
			ID:   arg.ID,
			Name: arg.Name,
		}

		img, err := s.image.Update(updateImg)
		if err != nil {
			return err
		}

		retPost.Image = img
		retPost.Tags = tags

		return nil
	})

	// Return updated post and error
	return retPost, err
}

// Deletes an image post from the database
func (s *PostSQL) Delete(ctx context.Context, id int64) error {
	err := s.execTx(ctx, func(s *PostSQL) error {
		var err error

		// Tags from an Image need to be removed before deletion
		err = s.imageTag.DeleteAllFromImage(id)
		if err != nil {
			return err
		}

		err = s.image.Delete(id)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// Associates an image with a slice of tags
func (s *PostSQL) tagImagePost(imgID int64, args []models.CreateTagParams) ([]models.Tag, error) {
	var tags []models.Tag
	for _, tagArg := range args {
		tag, err := s.tag.Create(tagArg)
		if err != nil {
			return tags, err
		}

		// Associate tags with image
		imgTag := models.ImageTag{
			ImageID: imgID,
			TagID:   tag.ID,
		}
		if err = s.imageTag.Create(imgTag); err != nil {
			return tags, err
		}

		tags = append(tags, tag)
	}
	return tags, nil
}
