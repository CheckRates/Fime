package postgres

import (
	"database/sql"
	"testing"

	"github.com/checkrates/Fime/util"
	"github.com/stretchr/testify/require"
)

func createTestImage(t *testing.T) Image {
	user := createTestUser(t)

	args := CreateImageParams{
		Name:    util.RandomString(5),
		URL:     "www." + util.RandomString(10) + ".com",
		OwnerID: user.ID,
	}

	img, err := dal.CreateImage(args)
	require.NoError(t, err)
	require.NotZero(t, img.ID)
	require.Equal(t, args.Name, img.Name)
	require.Equal(t, args.URL, img.URL)
	require.Equal(t, args.OwnerID, img.OwnerID)

	return img
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

	img2Args := UpdateImageParams{
		ID:   img.ID,
		Name: util.RandomString(6),
	}

	beforeImg := img
	img, err := dal.UpdateImage(img2Args)

	require.NoError(t, err)
	require.Equal(t, img.ID, beforeImg.ID)
	require.Equal(t, img.Name, img2Args.Name)
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

	listArgs := ListImageParams{
		Limit:  5,
		Offset: 5,
	}

	images, err := dal.Images(listArgs)
	require.NoError(t, err)

	for _, img := range images {
		require.NotEmpty(t, img)
	}
}
