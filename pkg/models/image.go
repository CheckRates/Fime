package models

import "time"

type Image struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	URL       string    `db:"url"`
	OwnerID   int64     `db:"owner"`
	CreatedAt time.Time `db:"createdat"`
}

// Parameters to list all images from a repository
type ListImagesParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

// Parameters to create a new image in the repository
type CreateImageParams struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	OwnerID int64  `json:"userID"`
}

// Parameters to update an image in the repository
type UpdateImageParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Parameters to list an user's images
type ListUserImagesParams struct {
	UserID int64 `json:"userID"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}
