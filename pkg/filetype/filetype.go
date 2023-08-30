package filetype

import (
	"encoding/base64"
	"fmt"

	"github.com/gabriel-vasile/mimetype"
)

type imageType string

const (
	PNG     imageType = "image/png"
	JPG     imageType = "image/jpg"
	JPEG    imageType = "image/jpeg"
	Invalid imageType = "invalid"
)

func getValidImageTypes() []imageType {
	return []imageType{PNG, JPG, JPEG}
}

// Checks if a base64 encoded file header is from a valid image type
func IsValid(fileheader string) (bool, error) {
	imageType, err := Get(fileheader)
	if err != nil {
		return false, err
	}

	if imageType != Invalid {
		return true, nil
	}
	return false, nil
}

// Retrieves the type of image indentified in a base64 encoded fileheader
func Get(fileheader string) (imageType, error) {
	byteHeader, err := base64.StdEncoding.DecodeString(fileheader)
	if err != nil {
		return Invalid, fmt.Errorf("failed to decode a base64 string: %w", err)
	}

	mime := mimetype.Detect(byteHeader)

	for _, imgType := range getValidImageTypes() {
		if mime.Is(string(imgType)) {
			return imgType, nil
		}
	}
	return Invalid, nil
}
