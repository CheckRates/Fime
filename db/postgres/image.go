package postgres

import (
	"fmt"

	"github.com/checkrates/Fime/fime"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// NewImageStore returns the access point to all the images of Fime
func NewImageStore(db *sqlx.DB) *ImageStore {
	return &ImageStore{
		DB: db,
	}
}

// ImageStore is the database access point to the images
type ImageStore struct {
	*sqlx.DB
}

// Image return image by id
func (s *ImageStore) Image(id uuid.UUID) (fime.Image, error) {
	var i fime.Image
	if err := s.Get(&i, `SELECT * FROM images WHERE id=$1 LIMIT 1`, id); err != nil {
		return fime.Image{}, fmt.Errorf("error retrieving image: %w", err)
	}
	return i, nil
}

// Images return all images
func (s *ImageStore) Images(limit int, offset int) ([]fime.Image, error) {
	var ii []fime.Image
	if err := s.Get(&ii, `SELECT * FROM images ORDER BY id LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return []fime.Image{}, fmt.Errorf("error retrieving images: %w", err)
	}
	return ii, nil
}

// CreateImage uploads a new image to the database
func (s *ImageStore) CreateImage(i *fime.Image) error {
	if err := s.Get(i, `INSERT INTO images VALUES ($1, $2) RETURNING *`, i.Name, i.URL); err != nil {
		return fmt.Errorf("error inserting new image: %w", err)
	}
	return nil
}

// UpdateImage updates an image
func (s *ImageStore) UpdateImage(i *fime.Image) error {
	if err := s.Get(i, `UPDATE images SET name = $1 RETURNING *`, i.Name); err != nil {
		return fmt.Errorf("error updating image: %w", err)
	}
	return nil
}

// DeleteImage deletes an image from the database
func (s *ImageStore) DeleteImage(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM images WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting image: %w", err)
	}
	return nil
}
