package authorization

import (
	"context"

	"github.com/pkg/errors"
	"local.com/go-clean-lambda/internal/logger"
	"local.com/go-clean-lambda/internal/sdk/account"
)

// provided by app side.
const (
	AuthIndexApp         int = iota // index: 0
	AuthIndexAppDummy               // index: 1
	AuthIndexAppDummyNew            // index: 2
	AuthIndexAppPing                // index: 3
)

type RoleDummyClient struct {
	imap map[string][]int
}

func NewRoleDummyClient() *RoleDummyClient {
	defaultIndexMap := map[string][]int{
		account.User0ID:    {AuthIndexApp},
		account.User01ID:   {AuthIndexApp, AuthIndexAppDummy},
		account.User012ID:  {AuthIndexApp, AuthIndexAppDummy, AuthIndexAppDummyNew},
		account.User03ID:   {AuthIndexApp, AuthIndexAppPing},
		account.User013ID:  {AuthIndexApp, AuthIndexAppDummy, AuthIndexAppPing},
		account.User0123ID: {AuthIndexApp, AuthIndexAppDummy, AuthIndexAppDummyNew, AuthIndexAppPing},
		account.UserRootID: {AuthIndexApp, AuthIndexAppDummy, AuthIndexAppDummyNew, AuthIndexAppPing},
	}
	return &RoleDummyClient{
		imap: defaultIndexMap,
	}
}

func (c *RoleDummyClient) GetPermissionBit(ctx context.Context, permissions []int, isRoot bool) (uint64, error) {
	if isRoot {
		return GenerateRootBit(), nil
	}
	bit, err := GenerateGrantedBit(permissions)
	if err != nil {
		return NonePermissionBit, errors.Wrapf(
			err,
			"failed to generate permission bit. permissions: %s",
			logger.Pretty(permissions))
	}
	return bit, nil
}

func (c *RoleDummyClient) ListGrantedIndices(ctx context.Context, userID string) ([]int, error) {
	if userID == "" {
		return nil, errors.WithStack(ErrInvalidUserID)
	}
	gi, ok := c.imap[userID]
	if !ok {
		return []int{}, nil
	}
	return gi, nil
}
