package repository

import (
	"context"
	"testing"

	"github.com/Entetry/userService/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	user = model.User{ //nolint:gochecknoglobals //Explanation user for test
		ID:           uuid.New(),
		Username:     "YungLean",
		Email:        "jonahtanlendroyer@proton.me",
		PasswordHash: "blablabla",
	}
)

func TestUser_Create_And_GetByID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		_, err := dbPool.Exec(ctx, "TRUNCATE table users")
		require.NoError(t, err)
	}()
	t.Log("Given the need to test create user.")
	id, err := userRepository.Create(ctx, user.Username, user.PasswordHash, user.Email)
	require.NoError(t, err, "tested create function error")
	one, err := userRepository.GetByID(ctx, id)
	require.NoError(t, err, "tested get function error")
	require.Equal(t, user.Email, one.Email)
}

func TestUser_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		_, err := dbPool.Exec(ctx, "TRUNCATE table users")
		require.NoError(t, err)
	}()
	t.Log("Given the need to test delete company.")
	id, err := userRepository.Create(ctx, user.Username, user.PasswordHash, user.Email)
	require.NoError(t, err, "tested create function error")
	err = userRepository.Delete(ctx, id)
	require.NoError(t, err, "delete function error")
	_, err = userRepository.GetByID(ctx, id)
	require.Error(t, ErrUserNotFound, err)
}

func TestUser_Create_And_GetByUsername(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		_, err := dbPool.Exec(ctx, "TRUNCATE table users")
		require.NoError(t, err)
	}()
	t.Log("Given the need to test create user.")
	_, err := userRepository.Create(ctx, user.Username, user.PasswordHash, user.Email)
	require.NoError(t, err, "tested create function error")
	one, err := userRepository.GetByUsername(ctx, user.Username)
	require.NoError(t, err, "tested get function error")
	require.Equal(t, user.Email, one.Email)
}
