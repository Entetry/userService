// Package service contains service structs
package service

import (
	"context"
	"errors"
	"github.com/Entetry/userService/internal/model"
	"github.com/Entetry/userService/internal/repository"
	"regexp"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// EmailRegex used for email checks
const EmailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

var (
	// ErrEmailNotValid Not valid email error
	ErrEmailNotValid = errors.New("email Not valid")
	// ErrUserNotFound not found user error
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailAlreadyExist email already exist err
	ErrEmailAlreadyExist = errors.New("email already exists")
	// ErrUsernameAlreadyExist username already exist err
	ErrUsernameAlreadyExist = errors.New("username already exists")
)

// UserRepository user repository interface
type UserRepository interface {
	Create(ctx context.Context, username, pwdHash, email string) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

// User service struct
type User struct {
	userRepository UserRepository
	emailRegex     *regexp.Regexp
}

// NewUserService creates new User service
func NewUserService(userRepository UserRepository) *User {
	regex := regexp.MustCompile(EmailRegex)
	return &User{
		userRepository: userRepository,
		emailRegex:     regex}
}

// GetByID GetByID return user by its id
func (u *User) GetByID(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	user, err := u.userRepository.GetByID(ctx, ID)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrUserNotFound
	} else if err != nil {
		log.Errorf("User / GetById error: \n %v", err)
		return nil, err
	}
	return user, nil
}

// GetByUsername return user by its username
func (u *User) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := u.userRepository.GetByUsername(ctx, username)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrUserNotFound
	} else if err != nil {
		log.Errorf("User / GetByUsername error: \n %v", err)
		return nil, err
	}
	return user, nil
}

// Create save user to db
func (u *User) Create(ctx context.Context, username, password, email string) (uuid.UUID, error) {
	lcEmail := strings.ToLower(email)
	if !u.isValidEmail(lcEmail) {
		return uuid.Nil, ErrEmailNotValid
	}
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Errorf("User / Create / Failed to create user:\n %v", err)
		return uuid.Nil, err
	}
	id, err := u.userRepository.Create(ctx, username, string(pwdHash), lcEmail)
	switch {
	case errors.Is(err, repository.ErrEmailAlreadyExist):
		return uuid.Nil, ErrEmailAlreadyExist
	case errors.Is(err, repository.ErrUsernameAlreadyExist):
		return uuid.Nil, ErrUsernameAlreadyExist
	case err != nil:
		log.Errorf("User / GetById error: \n %v", err)
		return uuid.Nil, err
	}
	return id, err
}

// Delete delete user from db
func (u *User) Delete(ctx context.Context, ID uuid.UUID) error {
	return u.userRepository.Delete(ctx, ID)
}

func (u *User) isValidEmail(email string) bool {
	return u.emailRegex.MatchString(email)
}
