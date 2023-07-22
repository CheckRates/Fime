package tag

import (
	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/checkrates/Fime/pkg/storage"
)

type tagService struct {
	repo storage.TagRepository
}

func NewTagService(tag storage.TagRepository) service.TagUsecase {
	return tagService{
		repo: tag,
	}
}

// Takes a size and the page number to provide a subset of tags
func (t tagService) GetMultiple(size, page int) ([]models.Tag, error) {
	arg := models.ListTagsParams{
		Limit:  size,
		Offset: (page - 1) * size,
	}

	tags, err := t.repo.GetMultiple(arg)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// Takes a size and the page number to provide tags used by an user in all their posts
func (t tagService) GetUserTags(userId int64, size, page int) ([]models.Tag, error) {
	arg := models.ListUserTagsParams{
		ID:     userId,
		Limit:  size,
		Offset: (page - 1) * size,
	}

	tags, err := t.repo.GetUserTags(arg)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
