package account

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/pkg/errors"
)

const (
	User0ID    string = "user0"
	User01ID   string = "user01"
	User012ID  string = "user012"
	User03ID   string = "user03"
	User013ID  string = "user013"
	User0123ID string = "user0123"
	UserRootID string = "root"
)

type UserDummyClient struct {
	umap map[string]*UserDto
	pmap map[string]string
}

func NewUserDummmyClient() *UserDummyClient {
	defaultUserMap := map[string]*UserDto{
		User0ID:    {UserID: User0ID, Name: User0ID},
		User01ID:   {UserID: User01ID, Name: User01ID},
		User012ID:  {UserID: User012ID, Name: User012ID},
		User03ID:   {UserID: User03ID, Name: User03ID},
		User013ID:  {UserID: User013ID, Name: User013ID},
		User0123ID: {UserID: User0123ID, Name: User0123ID},
		UserRootID: {UserID: UserRootID, Name: UserRootID},
	}
	defaultPasswordMap := map[string]string{
		User0ID:    encodePasssword(User0ID),
		User01ID:   encodePasssword(User01ID),
		User012ID:  encodePasssword(User012ID),
		User03ID:   encodePasssword(User03ID),
		User013ID:  encodePasssword(User013ID),
		User0123ID: encodePasssword(User0123ID),
		UserRootID: encodePasssword(UserRootID),
	}
	return &UserDummyClient{
		umap: defaultUserMap,
		pmap: defaultPasswordMap,
	}
}

func (c *UserDummyClient) GetUser(ctx context.Context, userID string) (*UserDto, error) {
	if userID == "" {
		return nil, errors.WithStack(ErrInvalidUserID)
	}
	return c.umap[userID], nil
}

func (c *UserDummyClient) RegisterUser(ctx context.Context, user *UserDto, password string) error {
	if user == nil || user.UserID == "" || user.Name == "" {
		return errors.WithStack(ErrInvalidUserInfo)
	}
	if password == "" {
		return errors.WithStack(ErrInvalidPassword)
	}
	c.umap[user.UserID] = user
	c.pmap[user.UserID] = encodePasssword(password)
	return nil
}

func (c *UserDummyClient) VerifyPassword(ctx context.Context, userID string, password string) (bool, error) {
	if userID == "" {
		return false, errors.WithStack(ErrInvalidUserID)
	}
	p1 := c.pmap[userID]
	p2 := encodePasssword(password)
	return p1 == p2, nil
}

// encodePasssword
//
//	@param password
//	@return string
func encodePasssword(password string) string {
	hashed := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hashed)
}
