package authentication

import (
	"context"
	"errors"
	"time"
)

type AuthContextKey string

const (
	UserContextKey AuthContextKey = "user_context"
	UserIDKey      AuthContextKey = "user_id"
	ExpireDuration time.Duration  = time.Minute * 30
)

var (
	ErrBadClient    error = errors.New("bad client")
	ErrInvalidJwt   error = errors.New("invalid jwt value")
	ErrInvalidClaim error = errors.New("invalid claim")
	ErrExpiredClaim error = errors.New("claim is expired")
	ErrBlockedClaim error = errors.New("claim is blocked")
)

type AuthJwtClient interface {
	// Issue
	//  @param ctx
	//  @param claim
	//  @return string
	//  @return error
	Issue(ctx context.Context, claim *AuthJwtClaim) (string, error)

	// Block
	//  @param ctx
	//  @param tokenStr
	//  @return error
	Block(ctx context.Context, tokenStr string) error

	// Verify
	//  @param ctx
	//  @param tokenStr
	//  @return *AuthJwtClaim
	//  @return error ErrBadClient, ErrInvalidJwt, ErrInvalidClaim, ErrExpiredClaim, ErrBlockedClaim and others
	Verify(ctx context.Context, tokenStr string) (*AuthJwtClaim, error)
}
