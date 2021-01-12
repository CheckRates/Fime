package postgres

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// Test data access layer
var dal *Store

func TestMain(m *testing.M) {
	conn, err := sqlx.Open("postgres", config.New().Database.ConnString)
	if err != nil {
		log.Fatal("TEST: Cannot connect to the database: ", err)
	}

	dal, err = NewStore(conn)
	os.Exit(m.Run())
}

/*
*	Image Post Tests
 */

func createTestImagePost(t *testing.T) ImagePostResult {
	var err error
	user := createTestUser(t)

	var tagsArgs []CreateTagParams
	for i := 0; i < 3; i++ {
		arg := CreateTagParams{
			Name: util.RandomString(4),
		}
		tagsArgs = append(tagsArgs, arg)
	}

	postArgs := MakePostParams{
		Name:   "IMG_2020",
		URL:    "www.coolImage.com",
		UserID: user.ID,
		Tags:   tagsArgs,
	}

	newPost, err := dal.MakePostTx(context.Background(), postArgs)

	require.NoError(t, err)
	require.NotZero(t, newPost.Image.ID)

	require.Equal(t, newPost.Image.Name, postArgs.Name)
	require.Equal(t, newPost.Image.URL, postArgs.URL)
	require.Equal(t, newPost.Image.OwnerID, postArgs.UserID)
	require.Equal(t, len(newPost.Tags), len(postArgs.Tags))
	for i := 0; i < len(newPost.Tags); i++ {
		require.Equal(t, newPost.Tags[i].Name, postArgs.Tags[i].Name)
	}

	return newPost
}

func TestMakePostTx(t *testing.T) {
	createTestImagePost(t)
}

func TestGetImagePost(t *testing.T) {

	imgPost := createTestImagePost(t)
	imgPost2, err := dal.GetPostTx(context.Background(), imgPost.Image.ID)

	require.NoError(t, err)
	require.NotEmpty(t, imgPost2)

	require.Equal(t, imgPost.Image.ID, imgPost2.Image.ID)
	require.Equal(t, imgPost.Image.Name, imgPost2.Image.Name)
	require.Equal(t, imgPost.Image.URL, imgPost.Image.URL)
	require.Equal(t, imgPost.Image.OwnerID, imgPost2.Image.OwnerID)
	require.Equal(t, imgPost.Tags, imgPost2.Tags)
}

func TestUpdateImagePost(t *testing.T) {
	imgPost := createTestImagePost(t)

	postArgs := UpdatePostParams{
		ID:   imgPost.Image.ID,
		Name: util.RandomString(6),
		Tags: []CreateTagParams{
			{
				Name: util.RandomString(4),
			},
			{
				Name: util.RandomString(4),
			},
		},
	}

	updatedPost, err := dal.UpdatePostTx(context.Background(), postArgs)

	require.NoError(t, err)
	require.Equal(t, updatedPost.Image.ID, imgPost.Image.ID)
	require.NotEqual(t, updatedPost.Image.Name, imgPost.Image.Name)
	require.NotEmpty(t, updatedPost.Tags)
}

func TestDeleteImagePost(t *testing.T) {
	imgPost := createTestImagePost(t)

	err := dal.DeletePostTx(context.Background(), imgPost.Image.ID)
	require.NoError(t, err)

	imgPost2, err := dal.GetPostTx(context.Background(), imgPost.Image.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, imgPost2)
}

func TestImagePostList(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestImagePost(t)
	}

	listArgs := ListParams{
		Limit:  5,
		Offset: 5,
	}

	posts, err := dal.ListPostTx(context.Background(), listArgs)
	require.NoError(t, err)

	for _, post := range posts {
		require.NotEmpty(t, post)
		require.NotEmpty(t, post.Image)
	}
}
