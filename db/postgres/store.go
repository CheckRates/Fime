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

func (store *Store) execTx(ctx context.Context, fn func(*Store) error) error {
	tx, err := store.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(store)
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
 * 	Post an Image Transaction
 */

// MakePostParams provides all data for creating a image post in Fime
type MakePostParams struct {
	Name   string            `json:"name"`
	URL    string            `json:"url"`
	UserID int64             `json:"ownerID"`
	Tags   []CreateTagParams `json:"tags"`
}

// MakePostResult is the result of creating a post
type MakePostResult struct {
	Image Image `json:"image"`
	Tags  []Tag `json:"tags"`
}

// MakePostTx creates a image post in the database
func (s *Store) MakePostTx(ctx context.Context, arg MakePostParams) (MakePostResult, error) {
	var retPost MakePostResult

	err := s.execTx(ctx, func(store *Store) error {
		var err error

		// Create Image
		imageArgs := CreateImageParams{
			Name:    arg.Name,
			URL:     arg.URL,
			OwnerID: arg.UserID,
		}
		img, err := store.CreateImage(imageArgs)
		if err != nil {
			return err
		}

		var tags []Tag
		for i := 0; i < len(arg.Tags); i++ {
			// Create Tag -> Tag that already exist are returned and not recreated
			tag, err := s.CreateTag(arg.Tags[i])
			if err != nil {
				return err
			}

			// Associate tags with image
			imgTag := ImageTag{
				ImageID: img.ID,
				TagID:   tag.ID,
			}
			if err = s.CreateImageTag(imgTag); err != nil {
				return err
			}

			tags = append(tags, tag)
		}

		// No errors -> proceed to return post
		retPost.Image = img
		retPost.Tags = tags
		return nil
	})

	// Return created post and error
	return retPost, err
}
