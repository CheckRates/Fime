package postgres

import (
	"github.com/checkrates/Fime/pkg/models"
	"github.com/jmoiron/sqlx"
)

type ImageTagSQL struct {
	*sqlx.DB
}

// Returns the access point to Fime's ImageTags
func NewImageTagRepository(db *sqlx.DB) *ImageTagSQL {
	return &ImageTagSQL{
		DB: db,
	}
}

// Creates a new associative entity between an image and tag
func (s *ImageTagSQL) Create(it models.ImageTag) error {
	if _, err := s.Exec(`INSERT INTO image_tags VALUES ($1, $2)`, it.ImageID, it.TagID); err != nil {
		return err
	}
	return nil
}

// Dissociate an image from a tag
func (s *ImageTagSQL) Delete(it models.ImageTag) error {
	_, err := s.Exec(`DELETE FROM image_tags WHERE image_id = $1 AND tag_id = $2`, it.ImageID, it.TagID)
	if err != nil {
		return err
	}
	return nil
}

// Dissociate all tags from a image
func (s *ImageTagSQL) DeleteAllFromImage(imgID int64) error {
	_, err := s.Exec(`DELETE FROM image_tags WHERE image_id = $1`, imgID)
	if err != nil {
		return err
	}
	return nil
}

// Returns all the tags of a specific image
func (s *ImageTagSQL) GetTagsByImageID(imgID int64) ([]models.Tag, error) {
	var tt []models.Tag
	statement :=
		`SELECT t.* FROM tags t 
			INNER JOIN image_tags it ON t.id = it.tag_id
			WHERE it.image_id = $1`

	if err := s.Select(&tt, statement, imgID); err != nil {
		return []models.Tag{}, err
	}

	return tt, nil
}

// Returns all the tags of a specific image
func (s *ImageTagSQL) GetImagesByTagID(tagID int64) ([]models.Image, error) {
	var ii []models.Image
	statement :=
		`SELECT i.* FROM images i 
			INNER JOIN image_tags it ON i.id = it.image_id
			WHERE it.tag_id = $1`

	if err := s.Select(&ii, statement, tagID); err != nil {
		return []models.Image{}, err
	}

	return ii, nil
}
