package authentication

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type AuthJwtDuummyClient struct {
	publicKeyParam  string
	privateKeyParam string
	ssmClient       ssmiface.SSMAPI
	blockMap        map[string]int64
}

// NewAuthJwtDummyClient
//
//	@param publicKeyParam
//	@param privateKeyParam
//	@param ssmClient
//	@return *AuthJwtDuummyClient
func NewAuthJwtDummyClient(
	publicKeyParam string,
	privateKeyParam string,
	ssmClient ssmiface.SSMAPI,
) *AuthJwtDuummyClient {
	return &AuthJwtDuummyClient{
		publicKeyParam:  publicKeyParam,
		privateKeyParam: privateKeyParam,
		ssmClient:       ssmClient,
		blockMap:        make(map[string]int64),
	}
}

func (c *AuthJwtDuummyClient) Issue(ctx context.Context, claim *AuthJwtClaim) (string, error) {
	if claim == nil {
		return "", errors.WithStack(ErrInvalidClaim)
	}
	privateKey, err := c.getSSMParam(ctx, c.privateKeyParam)
	if err != nil {
		return "", errors.Wrap(err, "failed to get private key")
	}
	rsakey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		rootErr := errors.New(err.Error())
		return "", errors.Wrap(rootErr, "failed to parse private key")
	}
	curr := time.Now()
	claim.IssuesAt = curr.UnixNano()
	claim.ExpiresAt = curr.Add(ExpireDuration).UnixNano()
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	jwt, err := t.SignedString(rsakey)
	if err != nil {
		rootErr := errors.New(err.Error())
		return "", errors.Wrap(rootErr, "failed to issue jwt")
	}
	return jwt, nil
}

func (c *AuthJwtDuummyClient) Verify(ctx context.Context, tokenStr string) (*AuthJwtClaim, error) {
	if len(tokenStr) == 0 {
		return nil, errors.WithStack(ErrInvalidJwt)
	}
	blocked, err := c.isBlocked(tokenStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check block status")
	}
	if blocked {
		return nil, errors.WithStack(ErrBlockedClaim)
	}
	claim, err := c.parseJwt(ctx, tokenStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prase jwt as claim")
	}
	if claim == nil || claim.User == nil {
		return nil, errors.WithStack(ErrInvalidClaim)
	}
	if claim.ExpiresAt <= time.Now().UnixNano() {
		return nil, errors.WithStack(ErrExpiredClaim)
	}
	return claim, nil
}

func (c *AuthJwtDuummyClient) Block(ctx context.Context, tokenStr string) error {
	if len(tokenStr) == 0 {
		return nil
	}
	claim, err := c.parseJwt(ctx, tokenStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse jwt as claim")
	}
	if claim == nil {
		return errors.Wrapf(ErrInvalidJwt, "token: %s", tokenStr)
	}
	curr := time.Now().UnixNano()
	if claim.ExpiresAt <= curr {
		// already expired, no need to block it
		return nil
	}
	c.blockMap[tokenStr] = claim.ExpiresAt
	return nil
}

// isBlocked
//
//	@receiver c
//	@param tokenStr
//	@return bool
//	@return error
func (c *AuthJwtDuummyClient) isBlocked(tokenStr string) (bool, error) {
	if tokenStr == "" {
		return false, errors.WithStack(ErrInvalidJwt)
	}
	expireAt, ok := c.blockMap[tokenStr]
	if !ok {
		return false, nil
	}
	curr := time.Now().UnixNano()
	return expireAt > curr, nil
}

// parseJwt
//
//	@receiver c
//	@param ctx
//	@param tokenStr
//	@return *AuthJwtClaim
//	@return error
func (c *AuthJwtDuummyClient) parseJwt(ctx context.Context, tokenStr string) (*AuthJwtClaim, error) {
	publicKey, err := c.getSSMParam(ctx, c.publicKeyParam)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get public key")
	}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&AuthJwtClaim{},
		func(t *jwt.Token) (interface{}, error) {
			//nolint:wrapcheck
			return jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		},
	)
	if err != nil {
		rootErr := errors.New(err.Error())
		return nil, errors.Wrapf(
			rootErr,
			"failed to parse token: %s, by public secret key: %s",
			tokenStr, c.publicKeyParam)
	}
	claim, ok := token.Claims.(*AuthJwtClaim)
	if !ok {
		return nil, errors.WithStack(ErrInvalidClaim)
	}
	return claim, nil
}

// getSSMParam
//
//	@receiver c
//	@param ctx
//	@param paramName
//	@return string
//	@return error
func (c *AuthJwtDuummyClient) getSSMParam(ctx context.Context, paramName string) (string, error) {
	if paramName == "" {
		return "", errors.Wrapf(ErrInvalidJwt, "invalid key param name: %s", paramName)
	}
	if c.ssmClient == nil {
		return "", errors.Wrap(ErrBadClient, "no ssm client found")
	}
	input := &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(false),
	}
	output, err := c.ssmClient.GetParameterWithContext(ctx, input)
	if err != nil {
		return "", errors.Wrap(ErrInvalidJwt, err.Error())
	}
	if output == nil || output.Parameter == nil || output.Parameter.Value == nil {
		return "", errors.Wrap(ErrInvalidJwt, "no ssm param found")
	}
	value := *output.Parameter.Value
	return value, nil
}
