package models

type Tag struct {
	ID   int64  `db:"id"`
	Name string `db:"tag"`
}

// Parameters to list users from a repository
type ListTagsParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Parameters to create a tag
type CreateTagParams struct {
	Name string `json:"tag"`
}

// Parameters to list all tags a user has used in the images
type ListUserTagsParams struct {
	ID     int64 `json:"id"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
