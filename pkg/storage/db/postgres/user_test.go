package postgres

import (
	"database/sql"
	"testing"

	"github.com/checkrates/Fime/pkg/models"
	"github.com/checkrates/Fime/util"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T) models.User {
	hashedPassword, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := models.CreateUserParams{
		Name:           util.RandomString(6),
		Email:          util.RandomString(7) + "@email.com",
		HashedPassword: hashedPassword,
	}

	user, err := user.Create(arg)
	require.NoError(t, err)
	require.NotZero(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createTestUser(t)
}

func TestGetUserByEmail(t *testing.T) {

	user1 := createTestUser(t)
	user2, err := user.FindByEmail(user1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Email, user2.Email)
}

func TestGetUser(t *testing.T) {

	user1 := createTestUser(t)
	user2, err := user.FindById(user1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Name, user2.Name)
}

func TestUpdateUser(t *testing.T) {
	user1 := createTestUser(t)

	updateArgs := models.UpdateUserParams{
		ID:   user1.ID,
		Name: util.RandomString(6),
	}

	beforeUser := user1
	user1, err := user.Update(updateArgs)

	require.NoError(t, err)

	require.Equal(t, user1.ID, beforeUser.ID)
	require.Equal(t, user1.Name, updateArgs.Name)
}

func TestDeleteUser(t *testing.T) {

	user1 := createTestUser(t)

	err := user.Delete(user1.ID)
	require.NoError(t, err)

	user2, err := user.FindById(user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestUser(t)
	}

	listArgs := models.ListUserParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := user.GetMultiple(listArgs)
	require.NoError(t, err)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
