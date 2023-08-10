package models

type RequestUploadParams struct {
	Filename   string `json:"filename"`
	UserID     int64  `json:"userId"`
	Fileheader string `json:"fileheader"`
}

type RequestUploadResponse struct {
	UploadId string `json:"uploadId"`
	ImageKey string `json:"imageKey"`
}
