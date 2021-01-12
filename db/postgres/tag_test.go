package postgres

import (
	"database/sql"
	"testing"

	"github.com/checkrates/Fime/util"
	"github.com/stretchr/testify/require"
)

func createTestTag(t *testing.T) Tag {
	args := CreateTagParams{
		Name: util.RandomString(4),
	}

	tag, err := dal.CreateTag(args)
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
	tag2Args := CreateTagParams{
		Name: tag1.Name,
	}

	tag2, err := dal.CreateTag(tag2Args)

	require.NoError(t, err)
	require.NotEmpty(t, tag2)

	require.Equal(t, tag1.ID, tag2.ID)
	require.Equal(t, tag1.Name, tag2.Name)
}

func TestGetTag(t *testing.T) {

	tag := createTestTag(t)
	tag2, err := dal.Tag(tag.ID)

	require.NoError(t, err)
	require.NotEmpty(t, tag2)

	require.Equal(t, tag.ID, tag2.ID)
	require.Equal(t, tag.Name, tag2.Name)
}

func TestDeleteTag(t *testing.T) {

	tag := createTestTag(t)

	err := dal.DeleteTag(tag.ID)
	require.NoError(t, err)

	tag2, err := dal.Tag(tag.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, tag2)
}

func TestListTag(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestTag(t)
	}

	listArgs := ListParams{
		Limit:  5,
		Offset: 5,
	}

	tags, err := dal.Tags(listArgs)
	require.NoError(t, err)

	for _, tag := range tags {
		require.NotEmpty(t, tag)
	}
}
