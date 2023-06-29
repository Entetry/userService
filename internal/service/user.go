// Package service contains service structs
package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"userService/internal/model"
	"userService/internal/repository"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const EmailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

var (
	ErrEmailNotValid = errors.New("email Not valid")
	ErrUserNotFound  = errors.New("user not found")
)

// User service struct
type User struct {
	userRepository repository.UserRepository
	emailRegex     *regexp.Regexp
}

// NewUserService creates new User service
func NewUserService(userRepository repository.UserRepository) *User {
	regex := regexp.MustCompile(EmailRegex)
	return &User{
		userRepository: userRepository,
		emailRegex:     regex}
}

// GetByUsername return user by its username
func (u *User) GetByID(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	user, err := u.userRepository.GetByID(ctx, ID)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrUserNotFound
	} else if err != nil {
		log.Errorf("User / GetById error: \n %v", err)
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

	return u.userRepository.Create(ctx, username, string(pwdHash), lcEmail)
}

func (u *User) Delete(ctx context.Context, ID uuid.UUID) error {
	return u.userRepository.Delete(ctx, ID)
}

func (u *User) isValidEmail(email string) bool {
	return u.emailRegex.Match([]byte(email))
}
