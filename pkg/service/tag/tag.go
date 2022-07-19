package tag

import (
	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/checkrates/Fime/pkg/storage"
)

type TagService struct {
	tag storage.TagRepository
}

func NewTagService(tag storage.TagRepository) service.TagUsecase {
	return TagService{
		tag: tag,
	}
}

// Takes a size and the page number to provide a subset of tags
func (t TagService) GetMultiple(size, page int64) ([]models.Tag, error) {
	arg := models.ListTagsParams{
		Limit:  size,
		Offset: (page - 1) * size,
	}

	tags, err := t.tag.GetMultiple(arg)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// Takes a size and the page number to provide tags used by an user in all their posts
func (t TagService) GetUserTags(userId, size, page int64) ([]models.Tag, error) {
	arg := models.ListUserTagsParams{
		ID:     userId,
		Limit:  size,
		Offset: (page - 1) * size,
	}

	tags, err := t.tag.GetUserTags(arg)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
