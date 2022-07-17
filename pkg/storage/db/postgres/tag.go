package postgres

import (
	"strings"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/jmoiron/sqlx"
)

type TagSQL struct {
	*sqlx.DB
}

// Returns the access point to Fime's tags
func NewTagRepository(db *sqlx.DB) *TagSQL {
	return &TagSQL{
		DB: db,
	}
}

// Returns tag by id, if found
func (s *TagSQL) FindById(id int64) (models.Tag, error) {
	var t models.Tag
	if err := s.Get(&t, `SELECT * FROM tags WHERE id=$1 LIMIT 1`, id); err != nil {
		return models.Tag{}, err
	}
	return t, nil
}

// Return a subset of tags
func (s *TagSQL) GetMultiple(args models.ListTagsParams) ([]models.Tag, error) {
	var tt []models.Tag
	if err := s.Select(&tt, `SELECT * FROM tags ORDER BY id LIMIT $1 OFFSET $2`, args.Limit, args.Offset); err != nil {
		return []models.Tag{}, err
	}
	return tt, nil
}

// Creates a new tag to the database, if does not exist. Returns newly or existing tag
func (s *TagSQL) Create(args models.CreateTagParams) (models.Tag, error) {
	statement :=
		`WITH s AS (
		SELECT * FROM tags WHERE tag=$1
	), 
	i as (
		INSERT INTO tags(tag) SELECT $1
			WHERE NOT EXISTS (SELECT 1 FROM s)
			RETURNING *
	)
	SELECT * FROM i UNION ALL SELECT * from s`

	var t models.Tag
	if err := s.Get(&t, statement, strings.ToLower(args.Name)); err != nil {
		return t, err
	}
	return t, nil
}

// Retrieves all tags that a user has attribute to any of their posts
func (s *TagSQL) GetUserTags(arg models.ListUserTagsParams) ([]models.Tag, error) {
	statement :=
		`SELECT DISTINCT t.id, t.tag 
			FROM image_tags as it 
			INNER JOIN tags as t   ON it.tag_id = t.id
			INNER JOIN images as i ON i.id = it.image_id
		WHERE i."owner" = $1 LIMIT $2 OFFSET $3;`

	var tt []models.Tag
	if err := s.Select(&tt, statement, arg.ID, arg.Limit, arg.Offset); err != nil {
		return tt, err
	}
	return tt, nil
}

// Deletes an tag from the database
func (s *TagSQL) Delete(id int64) error {
	if _, err := s.Exec(`DELETE FROM tags WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
