// Package handler provides grpc user api
package handler

import (
	"context"
	"errors"
	"github.com/Entetry/userService/internal/repository"
	"github.com/Entetry/userService/internal/service"
	"github.com/Entetry/userService/protocol/userService"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// User handler user struct
type User struct {
	userService.UnimplementedUserServiceServer
	userService *service.User
}

// NewUser creates new user handler
func NewUser(user *service.User) *User {
	return &User{userService: user}
}

// GetByID Retrieves user based on given ID
func (u *User) GetByID(ctx context.Context, request *userService.GetByIDRequest) (*userService.GetByIDResponse, error) {
	id, err := uuid.Parse(request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := u.userService.GetByID(ctx, id)
	if errors.Is(err, service.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userService.GetByIDResponse{
		Uuid:  user.ID.String(),
		Name:  user.Username,
		Email: user.Email,
	}, nil
}

// GetByUsername Retrieves user based on given username
func (u *User) GetByUsername(ctx context.Context, request *userService.GetByUsernameRequest) (*userService.GetByUsernameResponse, error) {
	user, err := u.userService.GetByUsername(ctx, request.GetUsername())
	if errors.Is(err, service.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userService.GetByUsernameResponse{
		Uuid:         user.ID.String(),
		Name:         user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

// Create save new user
func (u *User) Create(ctx context.Context, request *userService.CreateRequest) (*userService.CreateResponse, error) {
	id, err := u.userService.Create(ctx, request.Username, request.Password, request.Email)
	if errors.Is(err, service.ErrEmailNotValid) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	switch {
	case errors.Is(err, repository.ErrEmailAlreadyExist):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, repository.ErrUsernameAlreadyExist):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrEmailNotValid):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case err != nil:
		log.Errorf("User / GetById error: \n %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userService.CreateResponse{
		Uuid: id.String(),
	}, nil
}

// Delete company based on given ID
func (u *User) Delete(ctx context.Context, request *userService.DeleteRequest) (*userService.DeleteResponse, error) {
	id, err := uuid.Parse(request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = u.userService.Delete(ctx, id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userService.DeleteResponse{}, nil
}
