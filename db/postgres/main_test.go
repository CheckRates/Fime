package postgres

import (
	"context"
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

func TestMakePostTx(t *testing.T) {
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
}
