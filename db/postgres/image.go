package postgres

import (
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

// CreateImageParams provides all info to create a new image in the db
type CreateImageParams struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	OwnerID int64  `json:"userID"`
}

// UpdateImageParams provides all info to change a image's name in the db
type UpdateImageParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Image return image by id
func (s *ImageStore) Image(id int64) (Image, error) {
	var i Image
	if err := s.Get(&i, `SELECT * FROM images WHERE id=$1 LIMIT 1`, id); err != nil {
		return Image{}, err
	}
	return i, nil
}

// Images return all images
func (s *ImageStore) Images(args ListParams) ([]Image, error) {
	var ii []Image
	if err := s.Select(&ii, `SELECT * FROM images ORDER BY id LIMIT $1 OFFSET $2`, args.Limit, args.Offset); err != nil {
		return []Image{}, err
	}
	return ii, nil
}

// CreateImage uploads a new image to the database
func (s *ImageStore) CreateImage(args CreateImageParams) (Image, error) {
	var i Image
	if err := s.Get(&i, `INSERT INTO images (name, url, owner) VALUES ($1, $2, $3) RETURNING *`, args.Name, args.URL, args.OwnerID); err != nil {
		return i, err
	}
	return i, nil
}

// UpdateImage updates an image
func (s *ImageStore) UpdateImage(args UpdateImageParams) (Image, error) {
	var i Image
	if err := s.Get(&i, `UPDATE images SET name = $1 WHERE id=$2 RETURNING *`, args.Name, args.ID); err != nil {
		return i, err
	}
	return i, nil
}

// DeleteImage deletes an image from the database
func (s *ImageStore) DeleteImage(id int64) error {
	if _, err := s.Exec(`DELETE FROM images WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
