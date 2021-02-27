package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/checkrates/Fime/db/mock"
	"github.com/checkrates/Fime/db/postgres"
	"github.com/checkrates/Fime/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetImagePostAPI(t *testing.T) {
	imgPost := randomImagePost()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mock.NewMockStore(ctrl)

	store.EXPECT().GetPostTx(gomock.Any(), gomock.Eq(imgPost.Image.ID)).
		Times(1).
		Return(imgPost, nil)

	// Start test server
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	// Make request to the test server
	url := fmt.Sprintf("/image/%d", imgPost.Image.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, imgPost)

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, imgPost postgres.ImagePostResult) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var retImgPost postgres.ImagePostResult
	err = json.Unmarshal(data, &retImgPost)
	require.NoError(t, err)
	require.Equal(t, imgPost, retImgPost)
}

func randomImagePost() postgres.ImagePostResult {
	return postgres.ImagePostResult{
		Image: postgres.Image{
			ID:        util.RandomNumber(1000),
			Name:      util.RandomString(6),
			URL:       util.RandomString(10),
			OwnerID:   util.RandomNumber(10),
			CreatedAt: time.Now(),
		},
		Tags: []postgres.Tag{
			{
				ID:   util.RandomNumber(10),
				Name: util.RandomString(6),
			},
			{
				ID:   util.RandomNumber(10),
				Name: util.RandomString(6),
			},
		},
	}
}
