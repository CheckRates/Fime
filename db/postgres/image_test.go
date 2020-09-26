package postgres

import (
	"database/sql"
	"testing"

	"github.com/checkrates/Fime/fime"
	"github.com/checkrates/Fime/util"
	"github.com/stretchr/testify/require"
)

func createTestImage(t *testing.T) fime.Image {
	user := createTestUser(t)

	image := fime.Image{
		Name:    util.RandomString(5),
		URL:     "www." + util.RandomString(10) + ".com",
		OwnerID: user.ID,
	}

	err := dal.CreateImage(&image)
	require.NoError(t, err)
	require.NotZero(t, image.ID)

	return image
}

func TestCreateImage(t *testing.T) {
	createTestImage(t)
}

func TestGetImage(t *testing.T) {

	img1 := createTestImage(t)
	img2, err := dal.Image(img1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, img2)

	require.Equal(t, img1.ID, img2.ID)
	require.Equal(t, img1.Name, img2.Name)
}

func TestUpdateImage(t *testing.T) {

	img := createTestImage(t)

	beforeImg := img
	name := util.RandomString(6)
	img.Name = name
	err := dal.UpdateImage(&img)

	require.NoError(t, err)

	require.Equal(t, img.ID, beforeImg.ID)
	require.Equal(t, img.Name, name)
}

func TestDeleteImage(t *testing.T) {

	img1 := createTestImage(t)

	err := dal.DeleteImage(img1.ID)
	require.NoError(t, err)

	img2, err := dal.Image(img1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, img2)
}

func TestImageList(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestImage(t)
	}

	images, err := dal.Images(5, 5)
	require.NoError(t, err)

	for _, img := range images {
		require.NotEmpty(t, img)
	}
}
