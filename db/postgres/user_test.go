package postgres

import (
	"database/sql"
	"testing"

	"github.com/checkrates/Fime/fime"
	"github.com/checkrates/Fime/util"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T) fime.User {
	user := fime.User{
		Name:     util.RandomString(6),
		Email:    util.RandomString(7) + "@email.com",
		Password: util.RandomString(8),
	}

	err := dal.CreateUser(&user)
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	return user
}

func TestCreateUser(t *testing.T) {
	createTestUser(t)
}

func TestGetUser(t *testing.T) {

	user := createTestUser(t)
	user2, err := dal.User(user.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Name, user2.Name)
}

func TestUpdateUser(t *testing.T) {
	user := createTestUser(t)

	beforeUser := user
	name := util.RandomString(6)
	user.Name = name
	err := dal.UpdateUser(&user)

	require.NoError(t, err)

	require.Equal(t, user.ID, beforeUser.ID)
	require.Equal(t, user.Name, name)
}

func TestDeleteUser(t *testing.T) {

	user := createTestUser(t)

	err := dal.DeleteUser(user.ID)
	require.NoError(t, err)

	user2, err := dal.User(user.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestUser(t)
	}

	users, err := dal.Users(5, 5)
	require.NoError(t, err)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
