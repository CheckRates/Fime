package postgres

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/fime"
	"github.com/checkrates/Fime/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var dal *Store

func init() {
	var err error
	dal, err = NewStore(config.New().Database.ConnString)
	if err != nil {
		log.Fatal("error opening connecting to db: ", err)
	}
}

func TestMain(m *testing.M) {
	dal.User(1)
	os.Exit(m.Run())
}

func TestMakePostTx(t *testing.T) {
	user := createTestUser(t)
	var tags []fime.Tag
	for i := 0; i < 3; i++ {
		newTag := fime.Tag{
			Name: util.RandomString(4),
		}
		tags = append(tags, newTag)
	}

	postArgs := MakePostParams{
		Name:   "IMG_2020",
		URL:    "www.coolImage.com",
		UserID: user.ID,
		Tags:   tags,
	}

	newPost, err := dal.MakePostTx(context.Background(), postArgs)

	require.NoError(t, err)
	require.NotZero(t, newPost.Image.ID)

	require.Equal(t, newPost.Image.Name, postArgs.Name)
	require.Equal(t, newPost.Image.URL, postArgs.URL)
	require.Equal(t, newPost.Image.OwnerID, postArgs.UserID)
	require.Equal(t, newPost.Tags, postArgs.Tags)
}
