package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/google/uuid"
)

type s3Service struct {
	AwsConfig     aws.Config
	Region        string
	Bucket        string
	UseAccelerate bool

	RequestUploadRequestDuration time.Duration
}

func NewS3Service(region, bucket, accessId, secret string) (service.BucketUsecase, error) {
	config, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessId, secret, ""),
		),
	)
	if err != nil {
		return nil, err
	}

	return s3Service{
		AwsConfig:     config,
		Region:        region,
		Bucket:        bucket,
		UseAccelerate: true,
	}, nil
}

func (s s3Service) InitiateUpload(ctx context.Context, uploadParams models.InitiateUploadParams) (*models.InitiateUploadResponse, error) {
	client := s3.NewFromConfig(s.AwsConfig)

	imageId := uuid.New().String()
	imageKey := fmt.Sprintf("%d/%s/%s", uploadParams.UserID, imageId, uploadParams.Filename)
	expiresAt := time.Now().Add(s.RequestUploadRequestDuration)

	params := s3.CreateMultipartUploadInput{
		Bucket:  &s.Bucket,
		Key:     &imageKey,
		Expires: &expiresAt,
	}

	resp, err := client.CreateMultipartUpload(ctx, &params)
	if err != nil {
		return nil, err
	}

	return &models.InitiateUploadResponse{
		UploadId: *resp.UploadId,
		ImageKey: *resp.Key,
	}, nil
}

func (s s3Service) GeneratePresignURLs(ctx context.Context, uploadId string, numParts int32) ([]models.PresignedURL, error) {
	// TODO: not sure if I should share a initialize client in the service struct
	client := s3.NewPresignClient(s3.NewFromConfig(s.AwsConfig))

	var presignedURLs []models.PresignedURL

	// TODO: Make it multi threaded operation
	for i := int32(0); i < numParts; {
		resp, err := client.PresignUploadPart(ctx, &s3.UploadPartInput{
			Bucket:     &s.Bucket,
			PartNumber: i + 1,
			UploadId:   &uploadId,
		})

		// TODO: Make x number of retries
		if err == nil {
			presignedURLs = append(presignedURLs, models.PresignedURL{
				Header: resp.SignedHeader,
				URL:    resp.URL,
			})
			i++
		}
	}

	return presignedURLs, nil
}

func (s s3Service) CompleteUpload(ctx context.Context, uploadId string) (string, error) {
	client := s3.NewFromConfig(s.AwsConfig)

	params := s3.CompleteMultipartUploadInput{
		UploadId: &uploadId,
		Bucket:   &s.Bucket,
	}

	resp, err := client.CompleteMultipartUpload(ctx, &params)
	if err != nil {
		return "", err
	}

	return *resp.Location, err
}

func (s s3Service) Delete(ctx context.Context, key string) error {
	client := s3.NewFromConfig(s.AwsConfig)

	params := &s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	}

	_, err := client.DeleteObject(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (s s3Service) AbortUpload(ctx context.Context, imageKey string, uploadId string) error {
	client := s3.NewFromConfig(s.AwsConfig)

	params := &s3.AbortMultipartUploadInput{
		Bucket:   &s.Bucket,
		Key:      &imageKey,
		UploadId: &uploadId,
	}

	_, err := client.AbortMultipartUpload(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
