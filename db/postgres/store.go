package postgres

import (
	"fmt"

	"github.com/checkrates/Fime/fime"
	"github.com/jmoiron/sqlx"
)

// NewStore returns all the data access points of Fime
func NewStore(dataSource string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSource)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Return the wrapper of the DB
	return &Store{
		ImageStore: NewImageStore(db),
		UserStore:  NewUserStore(db),
		TagStore:   NewTagStore(db),
	}, nil
}

// Store contain all the data access points
type Store struct {
	fime.UserStore
	fime.ImageStore
	fime.TagStore
}
