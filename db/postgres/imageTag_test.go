package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestImageTags(t *testing.T) []ImageTag {
	var err error
	var imageTags []ImageTag

	img := createTestImage(t)
	for i := 0; i < 3; i++ {
		tag := createTestTag(t)
		it := ImageTag{
			ImageID: img.ID,
			TagID:   tag.ID,
		}
		err = dal.CreateImageTag(it)
		imageTags = append(imageTags, it)
	}

	tags, err := dal.GetTagsByImageID(img.ID)

	require.NoError(t, err)
	require.NotEmpty(t, tags)

	return imageTags
}

func TestCreateImageTags(t *testing.T) {
	createTestImageTags(t)
}

func TestDeleteImageTag(t *testing.T) {
	var err error
	imgTags := createTestImageTags(t)

	imgID := imgTags[0].ImageID
	oldTags, err := dal.GetTagsByImageID(imgID)

	// Disassociate the first tag from image
	dal.DeleteImageTag(imgTags[0])

	tags, err := dal.GetTagsByImageID(imgID)

	require.NoError(t, err)
	require.NotEqual(t, len(oldTags), len(tags))
	require.NotEqual(t, oldTags[0].Name, tags[0].Name)
}

func TestGetImagesByTagID(t *testing.T) {
	var err error
	tag := createTestTag(t)
	imgs := [2]Image{createTestImage(t), createTestImage(t)}

	imgTag := ImageTag{
		ImageID: imgs[0].ID,
		TagID:   tag.ID,
	}

	imgTag2 := ImageTag{
		ImageID: imgs[1].ID,
		TagID:   tag.ID,
	}

	err = dal.CreateImageTag(imgTag)
	err = dal.CreateImageTag(imgTag2)
	retImgs, err := dal.GetImagesByTagID(tag.ID)

	require.NoError(t, err)
	require.Equal(t, len(retImgs), len(imgs))
	require.ElementsMatch(t, imgs, retImgs)
}

func TestGetTagsByImage(t *testing.T) {
	var err error
	img := createTestImage(t)
	tag := createTestTag(t)
	it := ImageTag{
		ImageID: img.ID,
		TagID:   tag.ID,
	}

	err = dal.CreateImageTag(it)
	imgs, err := dal.GetImagesByTagID(tag.ID)

	require.NoError(t, err)
	require.NotEmpty(t, imgs)
}
