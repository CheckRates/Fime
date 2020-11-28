package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/checkrates/Fime/fime"
	"github.com/jmoiron/sqlx"
)

// NewStore returns all the data access points of Fime
func NewStore(db *sqlx.DB) (*Store, error) {
	// Test connection
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

// Store contain all the data access points to the tables of Fime
type Store struct {
	db *sqlx.DB
	fime.UserStore
	fime.ImageStore
	fime.TagStore
	fime.ImageTagStore
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
	Name   string     `json:"name"`
	URL    string     `json:"url"`
	UserID int64      `json:"ownerID"`
	Tags   []fime.Tag `json:"tags"`
}

// MakePostResult is the result of creating a post
type MakePostResult struct {
	Image fime.Image `json:"image"`
	Tags  []fime.Tag `json:"tags"`
}

// MakePostTx creates a image post in the database
func (store *Store) MakePostTx(ctx context.Context, arg MakePostParams) (MakePostResult, error) {
	var retPost MakePostResult

	err := store.execTx(ctx, func(store *Store) error {
		var err error

		// Create Image
		newImg := fime.Image{
			Name:    arg.Name,
			URL:     arg.URL,
			OwnerID: arg.UserID,
		}
		err = store.CreateImage(&newImg)
		if err != nil {
			return err
		}

		// Create Tags -> Tags that already exist are returned
		for i := 0; i < len(arg.Tags); i++ {
			if err = store.CreateTag(&arg.Tags[i]); err != nil {
				return err
			}

			// Associate tags with image
			imgTag := fime.ImageTag{
				ImageID: newImg.ID,
				TagID:   arg.Tags[i].ID,
			}
			if err = store.CreateImageTag(imgTag); err != nil {
				return err
			}
		}

		// No errors -> proceed to return post
		retPost.Image = newImg
		retPost.Tags = arg.Tags
		return nil
	})

	// Return created post and error
	return retPost, err
}
