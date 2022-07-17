package postgres

import (
	"testing"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/stretchr/testify/require"
)

func createTestImageTags(t *testing.T) []models.ImageTag {
	var err error
	var imageTags []models.ImageTag

	img := createTestImage(t)
	for i := 0; i < 3; i++ {
		tag := createTestTag(t)
		it := models.ImageTag{
			ImageID: img.ID,
			TagID:   tag.ID,
		}
		err = imageTag.Create(it)
		require.NoError(t, err)

		imageTags = append(imageTags, it)
	}

	tags, err := imageTag.GetTagsByImageID(img.ID)

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
	oldTags, err := imageTag.GetTagsByImageID(imgID)
	require.NoError(t, err)

	// Disassociate the first tag from image
	imageTag.Delete(imgTags[0])

	tags, err := imageTag.GetTagsByImageID(imgID)

	require.NoError(t, err)
	require.NotEqual(t, len(oldTags), len(tags))
	require.NotEqual(t, oldTags[0].Name, tags[0].Name)
}

func TestGetImagesByTagID(t *testing.T) {
	var err error
	tag := createTestTag(t)
	imgs := [2]models.Image{createTestImage(t), createTestImage(t)}

	imgTag := models.ImageTag{
		ImageID: imgs[0].ID,
		TagID:   tag.ID,
	}

	imgTag2 := models.ImageTag{
		ImageID: imgs[1].ID,
		TagID:   tag.ID,
	}

	err = imageTag.Create(imgTag)
	err = imageTag.Create(imgTag2)
	retImgs, err := imageTag.GetImagesByTagID(tag.ID)

	require.NoError(t, err)
	require.Equal(t, len(retImgs), len(imgs))
	require.ElementsMatch(t, imgs, retImgs)
}

func TestGetTagsByImage(t *testing.T) {
	var err error
	img := createTestImage(t)
	tag := createTestTag(t)
	it := models.ImageTag{
		ImageID: img.ID,
		TagID:   tag.ID,
	}

	err = imageTag.Create(it)
	imgs, err := imageTag.GetImagesByTagID(tag.ID)

	require.NoError(t, err)
	require.NotEmpty(t, imgs)
}

func TestDeleteAllImageTags(t *testing.T) {
	var err error
	img := createTestImage(t)
	for i := 0; i < 3; i++ {
		tag := createTestTag(t)
		it := models.ImageTag{
			ImageID: img.ID,
			TagID:   tag.ID,
		}

		err = imageTag.Create(it)
		require.NoError(t, err)
	}

	tags, err := imageTag.GetTagsByImageID(img.ID)
	require.NoError(t, err)

	err = imageTag.DeleteAllFromImage(img.ID)
	require.NoError(t, err)

	curTags, err := imageTag.GetTagsByImageID(img.ID)

	require.NoError(t, err)
	require.NotEmpty(t, tags)
	require.Empty(t, curTags)
}
