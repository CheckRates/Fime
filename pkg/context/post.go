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
func NewPostUsecase(db *sqlx.DB) service.PostUsecase {
	return post.NewPostService(
		postgres.NewPostRepository(db),
		bucket.NewFileBucket("./tempImages"),
	)
}

// Returns a new HTTP implementation of the postPort
func NewPostPort(db *sqlx.DB) http.PostPort {
	return http.NewPostApi(
		NewPostUsecase(db),
	)
}
