// Package handler provides grpc user api
package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"userService/internal/service"
	"userService/protocol/userService"
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
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "input arguments are invalid")
	}

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

// Create save new user
func (u *User) Create(ctx context.Context, request *userService.CreateRequest) (*userService.CreateResponse, error) {
	id, err := u.userService.Create(ctx, request.Username, request.Password, request.Email)
	if errors.Is(err, service.ErrEmailNotValid) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
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
