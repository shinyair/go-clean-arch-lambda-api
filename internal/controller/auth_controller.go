package controller

import (
	nativeerr "errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"local.com/go-clean-lambda/internal/logger"
	"local.com/go-clean-lambda/internal/sdk/account"
	"local.com/go-clean-lambda/internal/sdk/authentication"
	"local.com/go-clean-lambda/internal/sdk/authorization"
)

const (
	AuthIndexApp         int = iota // index: 0
	AuthIndexAppDummy               // index: 1
	AuthIndexAppDummyNew            // index: 2
	AuthIndexAppPing                // index: 3
)

var ErrInvalidUserIDOrPassword error = nativeerr.New("invalid user id or password")

// AuthController
// works as extends MuxControllerImpl.
type AuthController struct {
	*MuxControllerImpl
	jwtClient  authentication.AuthJwtClient
	roleClient authorization.RoleClient
	userClient account.UserClient
}

// NewAuthController
//
//	@param logMdf
//	@param authMdf
//	@param jwtClient
//	@param roleClient
//	@param userClient
//	@return *AuthController
func NewAuthController(
	logMdf mux.MiddlewareFunc,
	authMdf mux.MiddlewareFunc,
	jwtClient authentication.AuthJwtClient,
	roleClient authorization.RoleClient,
	userClient account.UserClient,
) *AuthController {
	c := &AuthController{
		MuxControllerImpl: NewMuxControllerImpl(
			"/auth",
			make(map[string]map[string]*MuxRouterHandler),
		),
		jwtClient:  jwtClient,
		userClient: userClient,
		roleClient: roleClient,
	}
	c.AddMuxRouter("/login", []string{
		http.MethodPost,
	}, []mux.MiddlewareFunc{
		logMdf,
	}, func(w http.ResponseWriter, r *http.Request) error {
		return c.login(w, r)
	})
	c.AddMuxRouter("/logout", []string{
		http.MethodPost,
	}, []mux.MiddlewareFunc{
		logMdf,
		authMdf,
	}, func(w http.ResponseWriter, r *http.Request) error {
		return c.logout(w, r)
	})
	c.AddMuxRouter("/register", []string{
		http.MethodPost,
	}, []mux.MiddlewareFunc{
		logMdf,
	}, func(w http.ResponseWriter, r *http.Request) error {
		return c.register(w, r)
	})
	return c
}

// login
//
// curl -X POST {host}/auth/login -d "userId={xxx}&password={xxx}"
//
//	@receiver c
//	@param w
//	@param r
//	@return error
func (c *AuthController) login(w http.ResponseWriter, r *http.Request) error {
	userID := r.FormValue("userId")
	password := r.FormValue("password")
	ctx := r.Context()
	verifed, err := c.userClient.VerifyPassword(ctx, userID, password)
	if err != nil {
		return errors.Wrapf(err, "failed to verify password. user_id: %s, password: %s", userID, password)
	}
	if !verifed {
		errMsg := ErrInvalidUserIDOrPassword.Error()
		logger.Info("%s. user_id: %s, password: %s", errMsg, userID, password)
		return c.WriteResponse(w, errMsg)
	}
	indices, err := c.roleClient.ListGrantedIndices(ctx, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to list granted permissions. user_id: %s", userID)
	}
	bit, err := c.roleClient.GetPermissionBit(ctx, indices, userID == account.UserRootID)
	if err != nil {
		return errors.Wrapf(err, "failed to grant permission bit. user_id: %s", userID)
	}
	user, err := c.userClient.GetUser(ctx, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to get user info. user_id: %s", userID)
	}
	jwt, err := c.jwtClient.Issue(ctx, &authentication.AuthJwtClaim{
		User: &authentication.UserContext{
			UserID:        userID,
			UserName:      user.Name,
			Locale:        "en",
			ZoneID:        "Asia/Tokyo",
			PermissionBit: bit,
		},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to issue jwt. user_id: %s", userID)
	}
	return c.WriteResponse(w, jwt)
}

// logout
//
// curl -X POST {host}/auth/logout -H "Authorization: Bearer {jwt}"
//
//	@receiver c
//	@param w
//	@param r
//	@return error
func (c *AuthController) logout(w http.ResponseWriter, r *http.Request) error {
	jwth := r.Header.Get(authentication.JwtHeader)
	jwt := jwth[len(authentication.JwtHeaderPrefix):]
	ctx := r.Context()
	err := c.jwtClient.Block(ctx, jwt)
	if err != nil {
		return errors.Wrapf(err, "failed to logout. jwt: %s", jwt)
	}
	return c.WriteResponse(w, "logout done")
}

// register
//
// curl -X POST {host}/auth/register -d "userId={xxx}&name={xxx}&password={xxx}"
//
//	@receiver c
//	@param w
//	@param r
//	@return error
func (c *AuthController) register(w http.ResponseWriter, r *http.Request) error {
	userID := r.FormValue("userId")
	userName := r.FormValue("name")
	password := r.FormValue("password")
	ctx := r.Context()
	user := &account.UserDto{
		UserID: userID,
		Name:   userName,
	}
	err := c.userClient.RegisterUser(ctx, user, password)
	if err != nil {
		return errors.Wrapf(err, "failed to register new user: %s", logger.Pretty(user))
	}
	indices, err := c.roleClient.ListGrantedIndices(ctx, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to list granted permissions. user_id: %s", userID)
	}
	bit, err := c.roleClient.GetPermissionBit(ctx, indices, userID == account.UserRootID)
	if err != nil {
		return errors.Wrapf(err, "failed to grant permission bit. user_id: %s", userID)
	}
	jwt, err := c.jwtClient.Issue(ctx, &authentication.AuthJwtClaim{
		User: &authentication.UserContext{
			UserID:        userID,
			Locale:        "en",
			ZoneID:        "Asia/Tokyo",
			PermissionBit: bit,
		},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to issue jwt. user_id: %s", userID)
	}
	return c.WriteResponse(w, jwt)
}
