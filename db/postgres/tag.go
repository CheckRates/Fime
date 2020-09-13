package postgres

import (
	"fmt"

	"github.com/checkrates/Fime/fime"
	"github.com/google/uuid"
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

// Tag return tag by id
func (s *TagStore) Tag(id uuid.UUID) (fime.Tag, error) {
	var t fime.Tag
	if err := s.Get(&t, `SELECT * FROM tags WHERE id=$1 LIMIT 1`, id); err != nil {
		return fime.Tag{}, fmt.Errorf("error retrieving tag: %w", err)
	}
	return t, nil
}

// Tags return all tags
func (s *TagStore) Tags(limit int, offset int) ([]fime.Tag, error) {
	var tt []fime.Tag
	if err := s.Get(&tt, `SELECT * FROM tags ORDER BY id LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return []fime.Tag{}, fmt.Errorf("error retrieving tags: %w", err)
	}
	return tt, nil
}

// CreateTag uploads a new tag to the database
func (s *TagStore) CreateTag(t *fime.Tag) error {
	if err := s.Get(t, `INSERT INTO tags VALUES ($1) RETURNING *`, t.Name); err != nil {
		return fmt.Errorf("error inserting new tag: %w", err)
	}
	return nil
}

// DeleteTag deletes an tag from the database
func (s *TagStore) DeleteTag(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM tags WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting tag: %w", err)
	}
	return nil
}
