package postgres

import (
	"database/sql"
	"testing"

	"github.com/checkrates/Fime/util"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T) User {
	newUser := CreateUserParams{
		Name:     util.RandomString(6),
		Email:    util.RandomString(7) + "@email.com",
		Password: util.RandomString(8),
	}

	user, err := dal.CreateUser(newUser)
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

	updateArgs := UpdateUserParams{
		ID:   user.ID,
		Name: util.RandomString(6),
	}

	beforeUser := user
	user, err := dal.UpdateUser(updateArgs)

	require.NoError(t, err)

	require.Equal(t, user.ID, beforeUser.ID)
	require.Equal(t, user.Name, updateArgs.Name)
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

	listArgs := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := dal.Users(listArgs)
	require.NoError(t, err)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
