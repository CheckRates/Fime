package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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
	//ext := strings.Split(strings.TrimLeft(img.EncodedImg, "data:image/"), ";")[0]

	img.EncodedImg = strings.Split(img.EncodedImg, ",")[1]
	imgBuffer, err := base64.StdEncoding.DecodeString(img.EncodedImg)
	if err != nil {
		return "", err
	}

	// TODO: Check image size
	//if(len(imgBuffer) > maxImageSize)
	urlPath := fmt.Sprintf("/images/%d/%s%s", img.UserID, bson.NewObjectId().Hex(), img.Name)

	config := config.New()
	u := s3manager.NewUploader(server.aws)
	_, err = u.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.S3.Bucket),
		Key:    aws.String(urlPath),
		Body:   bytes.NewBuffer(imgBuffer),
	})

	if err != nil {
		return "", err
	}

	// retutn the complete url
	urlPath = config.S3.Bucket + "s3.amazonaws.com" + urlPath
	return urlPath, nil
}
