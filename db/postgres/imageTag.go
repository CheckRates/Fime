package postgres

import (
	"github.com/checkrates/Fime/fime"
	"github.com/jmoiron/sqlx"
)

// NewImageTagStore returns the access point to ImageTags of Fime
func NewImageTagStore(db *sqlx.DB) *ImageTagStore {
	return &ImageTagStore{
		DB: db,
	}
}

// ImageTagStore is the database access point to ImageTags
type ImageTagStore struct {
	*sqlx.DB
}

// CreateImageTag creates a new associative entity between an image and tag
func (s *ImageTagStore) CreateImageTag(it fime.ImageTag) error {
	if _, err := s.Exec(`INSERT INTO image_tags VALUES ($1, $2)`, it.ImageID, it.TagID); err != nil {
		return err
	}
	return nil
}

// DeleteImageTag dissociate an image from a tag
func (s *ImageTagStore) DeleteImageTag(it fime.ImageTag) error {
	_, err := s.Exec(`DELETE FROM image_tags WHERE image_id = $1 AND tag_id = $2`, it.ImageID, it.TagID)
	if err != nil {
		return err
	}
	return nil
}

// GetTagsByImageID returns all the tags of a specific image
func (s *ImageTagStore) GetTagsByImageID(imgID int64) ([]fime.Tag, error) {
	var tt []fime.Tag
	statement :=
		`SELECT t.* FROM tags t 
			INNER JOIN image_tags it ON t.id = it.tag_id
			WHERE it.image_id = $1`

	if err := s.Select(&tt, statement, imgID); err != nil {
		return []fime.Tag{}, err
	}

	return tt, nil
}

// GetImagesByTagID returns all the tags of a specific image
func (s *ImageTagStore) GetImagesByTagID(tagID int64) ([]fime.Image, error) {
	var ii []fime.Image
	statement :=
		`SELECT i.* FROM images i 
			INNER JOIN image_tags it ON i.id = it.image_id
			WHERE it.tag_id = $1`

	if err := s.Select(&ii, statement, tagID); err != nil {
		return []fime.Image{}, err
	}

	return ii, nil
}
