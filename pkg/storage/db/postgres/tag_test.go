package postgres

import (
	"database/sql"
	"testing"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/util"
	"github.com/stretchr/testify/require"
)

func createTestTag(t *testing.T) models.Tag {
	args := models.CreateTagParams{
		Name: util.RandomString(4),
	}

	tag, err := tag.Create(args)
	require.NoError(t, err)
	require.Equal(t, args.Name, tag.Name)
	require.NotZero(t, tag.ID)
	require.Equal(t, tag.Name, args.Name)

	return tag
}

func TestCreateTag(t *testing.T) {
	createTestTag(t)
}

func TestCreateDuplicatedTag(t *testing.T) {
	tag1 := createTestTag(t)
	tag2Args := models.CreateTagParams{
		Name: tag1.Name,
	}

	tag2, err := tag.Create(tag2Args)

	require.NoError(t, err)
	require.NotEmpty(t, tag2)

	require.Equal(t, tag1.ID, tag2.ID)
	require.Equal(t, tag1.Name, tag2.Name)
}

func TestGetTag(t *testing.T) {

	tag1 := createTestTag(t)
	tag2, err := tag.FindById(tag1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, tag2)

	require.Equal(t, tag1.ID, tag2.ID)
	require.Equal(t, tag1.Name, tag2.Name)
}

func TestDeleteTag(t *testing.T) {

	tag1 := createTestTag(t)

	err := tag.Delete(tag1.ID)
	require.NoError(t, err)

	tag2, err := tag.FindById(tag1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, tag2)
}

func TestListTag(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestTag(t)
	}

	listArgs := models.ListTagsParams{
		Limit:  5,
		Offset: 5,
	}

	tags, err := tag.GetMultiple(listArgs)
	require.NoError(t, err)

	for _, tag := range tags {
		require.NotEmpty(t, tag)
	}
}
