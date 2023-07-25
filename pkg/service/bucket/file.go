package bucket

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/google/uuid"
)

type fileBucketService struct {
	folderPath string
}

func NewFileBucket(path string) service.BucketUsecase {
	return fileBucketService{
		folderPath: path,
	}
}

func (b fileBucketService) RequestUpload(uploadParams models.RequestUploadParams) (string, error) {
	// Since this is the file io implementation, the image will be
	// stored locally without any validation on the contents
	path, err := b.uploadImage(uploadParams.ImgData)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (b fileBucketService) Get(id string) (string, error) {
	path := filepath.Join(b.folderPath, id)

	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	return path, nil
}

func (b fileBucketService) Delete(url string) error {
	err := os.Remove(url)
	if err != nil {
		return err
	}

	return nil
}

func (b fileBucketService) uploadImage(encodedImgData string) (string, error) {
	imgData, err := base64.StdEncoding.DecodeString(encodedImgData)
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatalln(err)
	}

	id := uuid.New().String()
	imagePath := filepath.Join(b.folderPath, id)
	out, err := os.Create(imagePath)
	if err != nil {
		return "", nil
	}
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 1

	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		return "", err
	}
	return imagePath, nil
}
