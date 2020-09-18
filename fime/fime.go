package fime

// UserStore interface
type UserStore interface {
	User(id int64) (User, error)
	Users(limit int, offset int) ([]User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id int64) error
}

// ImageStore Interface
type ImageStore interface {
	Image(id int64) (Image, error)
	Images(limit int, offset int) ([]Image, error)
	CreateImage(i *Image) error
	UpdateImage(i *Image) error
	DeleteImage(id int64) error
}

// TagStore Interface
type TagStore interface {
	Tag(id int64) (Tag, error)
	Tags(limit int, offset int) ([]Tag, error)
	CreateTag(t *Tag) error
	DeleteTag(id int64) error
}
