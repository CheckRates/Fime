package context

import (
	"github.com/checkrates/Fime/pkg/http"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/checkrates/Fime/pkg/service/tag"
	"github.com/checkrates/Fime/pkg/storage/db/postgres"
	"github.com/jmoiron/sqlx"
)

// Returns a new instance of the tag usecase
func NewTagUsecase(db *sqlx.DB) service.TagUsecase {
	return tag.NewTagService(
		postgres.NewTagRepository(db),
	)
}

// Returns a new HTTP implementation of the tagPort
func NewTagPort(db *sqlx.DB) http.TagPort {
	return http.NewTagApi(
		NewTagUsecase(db),
	)
}
