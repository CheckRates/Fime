package fime

import "github.com/google/uuid"

// UserStore interface
type UserStore interface {
	User(id int64) (User, error)
	Users(limit int, offset int) ([]User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id uuid.UUID) error
}

// ImageStore Interface
type ImageStore interface {
	Image(id uuid.UUID) (Image, error)
	Images(limit int, offset int) ([]Image, error)
	CreateImage(i *Image) error
	UpdateImage(i *Image) error
	DeleteImage(id uuid.UUID) error
}

// TagStore Interface
type TagStore interface {
	Tag(id uuid.UUID) (Tag, error)
	Tags(limit int, offset int) ([]Tag, error)
	CreateTag(t *Tag) error
	DeleteTag(id uuid.UUID) error
}
