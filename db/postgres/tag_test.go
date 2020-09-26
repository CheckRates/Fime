package postgres

import (
	"database/sql"
	"testing"

	"github.com/checkrates/Fime/fime"
	"github.com/checkrates/Fime/util"
	"github.com/stretchr/testify/require"
)

func createTestTag(t *testing.T) fime.Tag {
	tag := fime.Tag{
		Name: util.RandomString(4),
	}

	err := dal.CreateTag(&tag)
	require.NoError(t, err)
	require.NotZero(t, tag.ID)

	return tag
}

func TestCreateTag(t *testing.T) {
	createTestTag(t)
}

func TestCreateDuplicatedTag(t *testing.T) {
	tag1 := createTestTag(t)
	tag2 := fime.Tag{
		Name: tag1.Name,
	}

	err := dal.CreateTag(&tag2)

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

	tag2, err := dal.User(tag.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, tag2)
}

func TestListTag(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestTag(t)
	}

	tags, err := dal.Tags(5, 5)
	require.NoError(t, err)

	for _, tag := range tags {
		require.NotEmpty(t, tag)
	}
}
