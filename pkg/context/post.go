package context

import (
	"github.com/checkrates/Fime/pkg/http"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/checkrates/Fime/pkg/service/bucket"
	"github.com/checkrates/Fime/pkg/service/post"
	"github.com/checkrates/Fime/pkg/storage/db/postgres"
	"github.com/jmoiron/sqlx"
)

// Returns a new instance of the post usecase
func NewPostUsecase(db *sqlx.DB, region, bucketName, accessId, secret string) (service.PostUsecase, error) {
	bucket, err := bucket.NewS3Service(region, bucketName, accessId, secret)
	if err != nil {
		return nil, err
	}

	return post.NewPostService(
		postgres.NewPostRepository(db),
		bucket,
	), nil
}

// Returns a new HTTP implementation of the postPort
func NewPostPort(db *sqlx.DB, region, bucket, accessId, secret string) (http.PostPort, error) {
	post, err := NewPostUsecase(db, region, bucket, accessId, secret)
	if err != nil {
		return nil, err
	}

	return http.NewPostApi(post), nil
}
