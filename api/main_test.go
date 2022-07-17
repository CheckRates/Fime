package api

import (
	"testing"
	"time"

	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/db/postgres"
	"github.com/checkrates/Fime/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store postgres.Store) *Server {
	config := config.Config{
		Token: config.JWT{
			AccessSecret:     util.RandomString(32),
			AccessExpiration: time.Minute,
		},
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
