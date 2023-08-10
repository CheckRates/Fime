package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/checkrates/Fime/pkg/filetype"
	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/google/uuid"
)

type s3Service struct {
	Region        string
	Bucket        string
	Access        string
	Secret        string
	UseAccelerate bool

	RequestUploadRequestDuration time.Duration
}

func NewS3Service(region, bucket, accessId, secret string) service.BucketUsecase {
	return s3Service{
		Region:        region,
		Bucket:        bucket,
		Access:        accessId,
		Secret:        secret,
		UseAccelerate: true,
	}
}

func (s s3Service) RequestUpload(uploadParams models.RequestUploadParams) (*models.RequestUploadResponse, error) {
	valid, err := filetype.IsValid(uploadParams.Fileheader)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file")
	}
	if !valid {
		return nil, fmt.Errorf("image file type not supported")
	}

	// This should definetely not be initiated per call
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Not sure if it should be initiated all the time for a call
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = s.Region
		o.UseAccelerate = s.UseAccelerate
		o.Credentials = aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(s.Access, s.Secret, ""))
	})

	imageId := uuid.New().String()
	imageKey := fmt.Sprintf("%d/%s/%s", uploadParams.UserID, imageId, uploadParams.Filename)
	expiresAt := time.Now().Add(s.RequestUploadRequestDuration)

	params := s3.CreateMultipartUploadInput{
		Bucket:  &s.Bucket,
		Key:     &imageKey,
		Expires: &expiresAt,
	}

	resp, err := client.CreateMultipartUpload(context.TODO(), &params)
	if err != nil {
		return &models.RequestUploadResponse{}, err
	}

	return &models.RequestUploadResponse{
		UploadId: *resp.UploadId,
		ImageKey: *resp.Key,
	}, nil
}

func (s s3Service) Get(id string) (string, error) {
	return "", nil
}

func (s s3Service) Delete(url string) error {
	return nil
}

func (s s3Service) getPresignedURLs(numFileParts int, uploadId, imageKey string) (string, error) {
	return "", nil
}
