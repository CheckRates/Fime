package postgres

import "time"

// User of the Fime app
type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time `db:"createdAt"`
}

// Image is a image
type Image struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	URL       string    `db:"url"`
	OwnerID   int64     `db:"owner"`
	CreatedAt time.Time `db:"createdat"`
}

// ImageTag associative entity
type ImageTag struct {
	ImageID int64
	TagID   int64
}

// Tag image tag
type Tag struct {
	ID   int64  `db:"id"`
	Name string `db:"tag"`
}
