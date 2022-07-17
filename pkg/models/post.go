package models

type ImagePost struct {
	Image Image `json:"image"`
	Tags  []Tag `json:"tags"`
}

// Provides all data for creating a image post in Fime
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
