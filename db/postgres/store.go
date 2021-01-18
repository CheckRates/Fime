package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Store contains all the data access points to the tables of Fime
type Store struct {
	db *sqlx.DB
	*UserStore
	*ImageStore
	*TagStore
	*ImageTagStore
}

// NewStore returns all the data access points of Fime
func NewStore(db *sqlx.DB) (*Store, error) {
	// Test db connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Return the wrapper of the DB
	return &Store{
		db:            db,
		ImageStore:    NewImageStore(db),
		UserStore:     NewUserStore(db),
		TagStore:      NewTagStore(db),
		ImageTagStore: NewImageTagStore(db),
	}, nil
}

func (s *Store) execTx(ctx context.Context, fn func(*Store) error) error {
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

/*
 * 	Image Post Transactions
 */

// MakePostParams provides all data for creating a image post in Fime
type MakePostParams struct {
	Name   string            `json:"name"`
	URL    string            `json:"url"`
	UserID int64             `json:"ownerID"`
	Tags   []CreateTagParams `json:"tags"`
}

// ImagePostResult is the result of creating a post
type ImagePostResult struct {
	Image Image `json:"image"`
	Tags  []Tag `json:"tags"`
}

// UpdatePostParams contains all the info to update an image post
type UpdatePostParams struct {
	ID   int64             `json:"id"`
	Name string            `json:"name"`
	Tags []CreateTagParams `json:"tags"`
}

// ListParams provides all the params to list objects from the db
type ListParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

// ListUserPostsParams provides all the params to list an user's posts
type ListUserPostsParams struct {
	UserID int64 `json:"userID"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

// MakePostTx creates a image post from the database
func (s *Store) MakePostTx(ctx context.Context, arg MakePostParams) (ImagePostResult, error) {
	var retPost ImagePostResult

	err := s.execTx(ctx, func(s *Store) error {
		var err error

		// Create Image
		imageArgs := CreateImageParams{
			Name:    arg.Name,
			URL:     arg.URL,
			OwnerID: arg.UserID,
		}
		img, err := s.CreateImage(imageArgs)
		if err != nil {
			return err
		}

		tags, err := s.tagImagePost(img.ID, arg.Tags)

		// No errors -> proceed to return post
		retPost.Image = img
		retPost.Tags = tags
		return nil
	})

	// Return created post and error
	return retPost, err
}

// GetPostTx gets a image post from the database
func (s *Store) GetPostTx(ctx context.Context, id int64) (ImagePostResult, error) {
	var retPost ImagePostResult

	err := s.execTx(ctx, func(s *Store) error {
		var err error

		img, err := s.Image(id)
		if err != nil {
			return err
		}

		tags, err := s.GetTagsByImageID(img.ID)
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

// ListPostTx gets a list of image posts from the database
func (s *Store) ListPostTx(ctx context.Context, arg ListParams) ([]ImagePostResult, error) {
	var retPost []ImagePostResult

	err := s.execTx(ctx, func(s *Store) error {
		var err error

		imgs, err := s.Images(arg)
		if err != nil {
			return err
		}

		// Get every image posts tags
		for _, img := range imgs {
			tags, err := s.GetTagsByImageID(img.ID)
			if err != nil {
				return err
			}

			post := ImagePostResult{
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

// ListUserPostTx gets a list of all the user's image posts from the database
func (s *Store) ListUserPostTx(ctx context.Context, arg ListUserPostsParams) ([]ImagePostResult, error) {
	var retPost []ImagePostResult

	err := s.execTx(ctx, func(s *Store) error {
		var err error

		imgs, err := s.ImagesByUser(ListUserImagesParams{
			UserID: arg.UserID,
			Offset: arg.Offset,
			Limit:  arg.Limit,
		})
		if err != nil {
			return err
		}

		// No images post case
		if len(imgs) == 0 {
			return nil
		}

		// Get every image posts tags
		for _, img := range imgs {
			tags, err := s.GetTagsByImageID(img.ID)
			if err != nil {
				return err
			}

			post := ImagePostResult{
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

// DeletePostTx deletes an image post from the database
func (s *Store) DeletePostTx(ctx context.Context, id int64) error {
	err := s.execTx(ctx, func(s *Store) error {
		var err error

		// Tags from an Image need to be removed before deletion
		err = s.DeleteAllImageTags(id)
		if err != nil {
			return err
		}

		err = s.DeleteImage(id)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// UpdatePostTx updates image post's tags info and image name
func (s *Store) UpdatePostTx(ctx context.Context, arg UpdatePostParams) (ImagePostResult, error) {
	var retPost ImagePostResult

	err := s.execTx(ctx, func(s *Store) error {
		var err error

		err = s.DeleteAllImageTags(arg.ID)
		if err != nil {
			return err
		}

		tags, err := s.tagImagePost(arg.ID, arg.Tags)
		if err != nil {
			return err
		}

		updateImg := UpdateImageParams{
			ID:   arg.ID,
			Name: arg.Name,
		}

		img, err := s.UpdateImage(updateImg)
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

// tagImagePost associates an image with tags
func (s *Store) tagImagePost(imgID int64, args []CreateTagParams) ([]Tag, error) {
	var tags []Tag
	for _, tagArg := range args {
		// Create Tag -> Tag that already exist are returned and not recreated
		tag, err := s.CreateTag(tagArg)
		if err != nil {
			return tags, err
		}

		// Associate tags with image
		imgTag := ImageTag{
			ImageID: imgID,
			TagID:   tag.ID,
		}
		if err = s.CreateImageTag(imgTag); err != nil {
			return tags, err
		}

		tags = append(tags, tag)
	}
	return tags, nil
}
