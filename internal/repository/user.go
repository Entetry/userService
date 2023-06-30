// Package repository contains postgres operations for user
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Entetry/userService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// constraintViolation error code
const constraintViolation = "23505"

var (
	// ErrUserNotFound tells that user is not found
	ErrUserNotFound = errors.New("user Not found")
	// ErrEmailAlreadyExist tells that email already exist
	ErrEmailAlreadyExist = errors.New("email already exists")
	// ErrUsernameAlreadyExist tells that email already exist
	ErrUsernameAlreadyExist = errors.New("username already exists")
)

// User User postgres repository struct
type User struct {
	db *pgxpool.Pool
}

// NewUserRepository creates new user repository object
func NewUserRepository(db *pgxpool.Pool) *User {
	return &User{
		db: db,
	}
}

// Create insert user record in db
func (u *User) Create(ctx context.Context, username, pwdHash, email string) (uuid.UUID, error) {
	var user model.User
	user.ID = uuid.New()
	user.PasswordHash = pwdHash
	user.Email = email
	user.Username = username
	_, err := u.db.Exec(ctx, `INSERT INTO users (id, username, email, passwordHash) VALUES ($1, $2, $3, $4)`,
		user.ID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		pqErr, ok := err.(*pgconn.PgError)
		if ok && pqErr.Code == constraintViolation {
			if pqErr.ConstraintName == "email_unique" {
				return uuid.Nil, ErrEmailAlreadyExist
			} else if pqErr.ConstraintName == "username_unique" {
				return uuid.Nil, ErrUsernameAlreadyExist
			}
		}

		return uuid.Nil, fmt.Errorf("cannot create User: %v", err)
	}
	return user.ID, nil
}

// GetByID return user by its id
func (u *User) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := u.db.QueryRow(ctx,
		`SELECT id, username, email, passwordHash FROM users WHERE id = $1`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error in GetByID: %v", err)
	}
	return &user, nil
}

// GetByUsername return user by its username
func (u *User) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := u.db.QueryRow(ctx, `SELECT id, username, email, passwordHash FROM users WHERE username = $1`, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("can't GetByUsername: %v", err)
	}
	return &user, nil
}

// Delete delete user by its id
func (u *User) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := u.db.Exec(ctx, "DELETE FROM users WHERE id =$1", id)
	if err != nil {
		return fmt.Errorf("cannot delete User with id %s: %v", id, err)
	}
	return nil
}
