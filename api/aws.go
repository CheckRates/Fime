package api

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/checkrates/Fime/config"
	"github.com/globalsign/mgo/bson"
)

const maxImageSize = 1000000000

func connectAWS() (*session.Session, error) {
	config := config.New()
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.S3.Region),
		Credentials: credentials.NewStaticCredentials(config.S3.Access, config.S3.Secret, ""),
	})

	if err != nil {
		return s, err
	}

	return s, nil
}

// UploadImage takes the base64 image content from a image post request
// and uploads to a S3 Bucket
func (server *Server) UploadImage(img *postImageParams) (string, error) {
	// Gets the file format extension (e.g png, jpeg)
	ext := strings.Split(strings.TrimLeft(img.EncodedImg, "data:image/"), ";")[0]

	img.EncodedImg = strings.Split(img.EncodedImg, ",")[1]
	imgBuffer, err := base64.StdEncoding.DecodeString(img.EncodedImg)
	if err != nil {
		return "", err
	}

	// TODO: Check image size
	//if(len(imgBuffer) > maxImageSize)

	imageKey := bson.NewObjectId().Hex()
	urlPath := fmt.Sprintf("/images/%d/%s.%s", img.UserID, imageKey, ext)

	config := config.New()
	u := s3manager.NewUploader(server.aws)
	_, err = u.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.S3.Bucket),
		Key:    aws.String(urlPath),
		Body:   bytes.NewBuffer(imgBuffer),
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		return "", err
	}

	// Make the full path to the S3 bucket
	urlPath = "https://" + config.S3.Bucket + ".s3.amazonaws.com" + urlPath
	return urlPath, nil
}

// DownloadImage takes the image path (key) and request to the S3 Bucket
// and returns an image
func (server *Server) DownloadImage(url string) (string, error) {
	img, err := os.Create(url)
	if err != nil {
		return "", err
	}

	config := config.New()
	d := s3manager.NewDownloader(server.aws)
	_, err = d.Download(img, &s3.GetObjectInput{
		Bucket: aws.String(config.S3.Bucket),
		Key:    aws.String(url),
	})

	if err != nil {
		return "", err
	}

	// Read image into bytes and encode it
	r := bufio.NewReader(img)
	buf, _ := ioutil.ReadAll(r)
	encodedImg := base64.StdEncoding.EncodeToString(buf)

	return encodedImg, nil
}

// DeleteImage takes an id and delete an image from the S3 bucket
func (server *Server) DeleteImage(id int64) error {
	svc := s3.New(server.aws)

	// Get post url
	img, err := server.store.Image(id)
	if err != nil {
		return err
	}

	key := strings.Split(img.URL, ".com")[1]
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(config.New().S3.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	return nil
}
