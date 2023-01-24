package account

import (
	"context"
	"errors"
)

var (
	ErrInvalidUserID   error = errors.New("invalid user id")
	ErrInvalidUserInfo error = errors.New("invalid user info")
	ErrInvalidPassword error = errors.New("invalid user password")
)

type UserDto struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

type UserClient interface {
	// GetUser
	//  @param ctx
	//  @param userID
	//  @return *UserDto
	//  @return error ErrInvalidUserID and others
	GetUser(ctx context.Context, userID string) (*UserDto, error)

	// RegisterUser
	//  @param ctx
	//  @param user
	//  @param password
	//  @return error ErrInvalidUserInfo, ErrInvalidPassword and others
	RegisterUser(ctx context.Context, user *UserDto, password string) error

	// VerifyPassword
	//  @param ctx
	//  @param userID
	//  @param password
	//  @return bool
	//  @return error ErrInvalidUserID and others
	VerifyPassword(ctx context.Context, userID string, password string) (bool, error)
}
