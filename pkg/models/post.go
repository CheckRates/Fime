package models

type ImagePost struct {
	Image Image `json:"image"`
	Tags  []Tag `json:"tags"`
}

type RequestUploadParams struct {
	Filename         string `json:"filename"`
	EncodedImgHeader string `json:"imageHeaderBase64"`
	ImageSize        int64  `json:"imageSize"`
	UserId           int64  `json:"userId"`
}

type RequestUploadResponse struct {
	Filename string `json:"filename"`
	NumParts int    `json:"numParts"`
	UploadId string `json:"uploadId"`
	ImageKey string `json:"imageKey"`
}

type CompleteUploadParams struct {
	Filename string            `json:"filename"`
	UploadID string            `json:"uploadId"`
	UserID   int64             `json:"ownerID"`
	Tags     []CreateTagParams `json:"tags"`
}

type CreatePostParams struct {
	Filename string            `json:"filename"`
	URL      string            `json:"url"`
	UserID   int64             `json:"ownerID"`
	Tags     []CreateTagParams `json:"tags"`
}

type UpdatePostParams struct {
	ID   int64             `json:"id"`
	Name string            `json:"name"`
	Tags []CreateTagParams `json:"tags"`
}
