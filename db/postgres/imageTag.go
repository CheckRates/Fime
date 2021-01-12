package postgres

import (
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
func (s *ImageTagStore) CreateImageTag(it ImageTag) error {
	if _, err := s.Exec(`INSERT INTO image_tags VALUES ($1, $2)`, it.ImageID, it.TagID); err != nil {
		return err
	}
	return nil
}

// DeleteImageTag dissociate an image from a tag
func (s *ImageTagStore) DeleteImageTag(it ImageTag) error {
	_, err := s.Exec(`DELETE FROM image_tags WHERE image_id = $1 AND tag_id = $2`, it.ImageID, it.TagID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAllImageTags dissociate all tags from a image
func (s *ImageTagStore) DeleteAllImageTags(imgID int64) error {
	_, err := s.Exec(`DELETE FROM image_tags WHERE image_id = $1`, imgID)
	if err != nil {
		return err
	}
	return nil
}

// GetTagsByImageID returns all the tags of a specific image
func (s *ImageTagStore) GetTagsByImageID(imgID int64) ([]Tag, error) {
	var tt []Tag
	statement :=
		`SELECT t.* FROM tags t 
			INNER JOIN image_tags it ON t.id = it.tag_id
			WHERE it.image_id = $1`

	if err := s.Select(&tt, statement, imgID); err != nil {
		return []Tag{}, err
	}

	return tt, nil
}

// GetImagesByTagID returns all the tags of a specific image
func (s *ImageTagStore) GetImagesByTagID(tagID int64) ([]Image, error) {
	var ii []Image
	statement :=
		`SELECT i.* FROM images i 
			INNER JOIN image_tags it ON i.id = it.image_id
			WHERE it.tag_id = $1`

	if err := s.Select(&ii, statement, tagID); err != nil {
		return []Image{}, err
	}

	return ii, nil
}
