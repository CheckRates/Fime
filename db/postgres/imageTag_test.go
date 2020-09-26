package postgres

import (
	"testing"

	"github.com/checkrates/Fime/fime"
	"github.com/stretchr/testify/require"
)

func TestCreateImageTags(t *testing.T) {
	var err error
	img := createTestImage(t)
	for i := 0; i < 3; i++ {
		tag := createTestTag(t)
		it := fime.ImageTag{
			ImageID: img.ID,
			TagID:   tag.ID,
		}
		err = dal.CreateImageTag(it)
	}

	tags, err := dal.GetTagsByImageID(img.ID)

	require.NoError(t, err)
	require.NotEmpty(t, tags)
}

func TestGetImageByTagID(t *testing.T) {
	var err error
	img := createTestImage(t)
	tag := createTestTag(t)
	it := fime.ImageTag{
		ImageID: img.ID,
		TagID:   tag.ID,
	}

	err = dal.CreateImageTag(it)
	imgs, err := dal.GetImagesByTagID(tag.ID)

	require.NoError(t, err)
	require.NotEmpty(t, imgs)
}
