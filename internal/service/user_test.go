package service

import (
	"context"
	"testing"

	"github.com/Entetry/userService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_Create_InvalidEmail(t *testing.T) {
	mockUsername := "test_user"
	mockPassword := "test_password"
	mockEmail := "invalid_email"
	mockUserRepository := mocks.NewUserRepository(t)
	userService := NewUserService(mockUserRepository)

	userID, err := userService.Create(context.Background(), mockUsername, mockPassword, mockEmail)
	assert.Equal(t, uuid.Nil, userID, "Expected empty user ID")
	assert.Equal(t, ErrEmailNotValid, err, "Expected ErrEmailNotValid error")
}
