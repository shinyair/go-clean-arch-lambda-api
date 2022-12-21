package authorization

import (
	"context"

	"github.com/pkg/errors"
)

var ErrInvalidUserID error = errors.New("invalid user id")

type RoleClient interface {
	// GetPermissionBit
	//  @param ctx
	//  @param indices
	//  @param isRoot
	//  @return uint64
	//  @return error
	GetPermissionBit(ctx context.Context, indices []int, isRoot bool) (uint64, error)

	// ListGrantedIndices
	//  @param ctx
	//  @param userID
	//  @return []int
	//  @return error ErrInvalidUserID and others
	ListGrantedIndices(ctx context.Context, userID string) ([]int, error)
}
