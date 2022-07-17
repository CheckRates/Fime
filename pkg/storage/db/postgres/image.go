package postgres

import (
	"github.com/checkrates/Fime/pkg/models"
	"github.com/jmoiron/sqlx"
)

type ImageSQL struct {
	*sqlx.DB
}

// Returns the access point to Fime's images
func NewImageRepository(db *sqlx.DB) *ImageSQL {
	return &ImageSQL{
		DB: db,
	}
}

// Return an image by id, if found
func (s *ImageSQL) FindById(id int64) (models.Image, error) {
	var i models.Image
	if err := s.Get(&i, `SELECT * FROM images WHERE id=$1 LIMIT 1`, id); err != nil {
		return models.Image{}, err
	}
	return i, nil
}

// Return a subset of images ordered by time created
func (s *ImageSQL) GetMultiple(args models.ListImagesParams) ([]models.Image, error) {
	var ii []models.Image
	if err := s.Select(&ii, `SELECT * FROM images ORDER BY createdat LIMIT $1 OFFSET $2`, args.Limit, args.Offset); err != nil {
		return []models.Image{}, err
	}
	return ii, nil
}

// Return all images from user ordered by time created
func (s *ImageSQL) GetByUser(args models.ListUserImagesParams) ([]models.Image, error) {
	var ii []models.Image
	err := s.Select(&ii, `SELECT * FROM images WHERE owner=$1 ORDER BY createdat LIMIT $2 OFFSET $3`,
		args.UserID, args.Limit, args.Offset)

	if err != nil {
		return []models.Image{}, err
	}
	return ii, nil
}

// Creates  a new image to the database
func (s *ImageSQL) Create(args models.CreateImageParams) (models.Image, error) {
	var i models.Image
	if err := s.Get(&i, `INSERT INTO images (name, url, owner) VALUES ($1, $2, $3) RETURNING *`, args.Name, args.URL, args.OwnerID); err != nil {
		return i, err
	}
	return i, nil
}

// Updates info of a image
func (s *ImageSQL) Update(args models.UpdateImageParams) (models.Image, error) {
	var i models.Image
	if err := s.Get(&i, `UPDATE images SET name = $1 WHERE id=$2 RETURNING *`, args.Name, args.ID); err != nil {
		return i, err
	}
	return i, nil
}

// Deletes an image from the database
func (s *ImageSQL) Delete(id int64) error {
	if _, err := s.Exec(`DELETE FROM images WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
