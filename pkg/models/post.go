package models

type ImagePost struct {
	Image Image `json:"image"`
	Tags  []Tag `json:"tags"`
}

// Provides all data to create a post in a service/usecase
type PostData struct {
	Name       string
	EncodedImg string
	UserId     int64
	Tags       []CreateTagParams
}

// Provides all data for saving a image post in a repository
type CreatePostParams struct {
	Name   string            `json:"name"`
	URL    string            `json:"url"`
	UserID int64             `json:"ownerID"`
	Tags   []CreateTagParams `json:"tags"`
}

type UpdatePostParams struct {
	ID   int64             `json:"id"`
	Name string            `json:"name"`
	Tags []CreateTagParams `json:"tags"`
}
