package postgres

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

// NewTagStore returns the access point to all the tags of Fime
func NewTagStore(db *sqlx.DB) *TagStore {
	return &TagStore{
		DB: db,
	}
}

// TagStore is the database access point to the tags
type TagStore struct {
	*sqlx.DB
}

// CreateTagParams provides all the info to create a tag
type CreateTagParams struct {
	Name string `json:"tag"`
}

// ListUserTagsParams contains all info to list all user params
type ListUserTagsParams struct {
	ID     int64 `json:"id"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

// Tag return tag by id
func (s *TagStore) Tag(id int64) (Tag, error) {
	var t Tag
	if err := s.Get(&t, `SELECT * FROM tags WHERE id=$1 LIMIT 1`, id); err != nil {
		return Tag{}, err
	}
	return t, nil
}

// Tags return all tags
func (s *TagStore) Tags(args ListParams) ([]Tag, error) {
	var tt []Tag
	if err := s.Select(&tt, `SELECT * FROM tags ORDER BY id LIMIT $1 OFFSET $2`, args.Limit, args.Offset); err != nil {
		return []Tag{}, err
	}
	return tt, nil
}

// CreateTag uploads a new tag to the database
func (s *TagStore) CreateTag(args CreateTagParams) (Tag, error) {
	// Insert a tag if it does NOT exist, Else return the existing tag
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

	var t Tag
	if err := s.Get(&t, statement, strings.ToLower(args.Name)); err != nil {
		return t, err
	}
	return t, nil
}

// GetUserTags retrieves all tags that a user has attribute to any of their posts
func (s *TagStore) GetUserTags(arg ListUserTagsParams) ([]Tag, error) {
	statement :=
		`SELECT DISTINCT t.id, t.tag 
			FROM image_tags as it 
			INNER JOIN tags as t   ON it.tag_id = t.id
			INNER JOIN images as i ON i.id = it.image_id
		WHERE i."owner" = $1 LIMIT $2 OFFSET $3;`

	var tt []Tag
	if err := s.Select(&tt, statement, arg.ID, arg.Limit, arg.Offset); err != nil {
		return tt, err
	}
	return tt, nil
}

// DeleteTag deletes an tag from the database
func (s *TagStore) DeleteTag(id int64) error {
	if _, err := s.Exec(`DELETE FROM tags WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
