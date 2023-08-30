package models

import "net/http"

type InitiateUploadParams struct {
	Filename   string
	UserID     int64
	Fileheader string
}

type InitiateUploadResponse struct {
	UploadId string
	ImageKey string
}

type PresignedURL struct {
	Header http.Header `json:"header"`
	URL    string      `json:"URL"`
}
