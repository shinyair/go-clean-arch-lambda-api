package controller

import (
	"context"
	nativeerr "errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"local.com/go-clean-lambda/internal/logger"
	"local.com/go-clean-lambda/internal/sdk/authentication"
	"local.com/go-clean-lambda/internal/sdk/authorization"
)

var (
	ErrInvalidAuthenticationHeader error = nativeerr.New("invalid authorization header")
	ErrUnauthenticated             error = nativeerr.New("unauthenticated")
	ErrForbidden                   error = nativeerr.New("forbidden")
	ErrUserContextNotFound         error = nativeerr.New("user context not found")
)

// GetLogMiddleware
//
//	@return mux.MiddlewareFunc
func GetLogMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtHeader := r.Header.Get(authentication.JwtHeader)
			logger.Info(
				"handle request. path: %s, method: %s, middleware: log, jwt header: %s",
				r.URL.Path, r.Method, jwtHeader)
			next.ServeHTTP(w, r)
		})
	}
}

// GetLoginAccessMiddleware
//
//	@param authClient
//	@return mux.MiddlewareFunc
func GetLoginAccessMiddleware(authClient authentication.AuthJwtClient) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// verify
			jwtHeader := r.Header.Get(authentication.JwtHeader)
			if !strings.HasPrefix(jwtHeader, authentication.JwtHeaderPrefix) {
				err := errors.WithStack(ErrInvalidAuthenticationHeader)
				logger.Error("failed to check auth. header: %s", err, jwtHeader)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			claim, err := authClient.Verify(r.Context(), jwtHeader[len(authentication.JwtHeaderPrefix):])
			if err != nil {
				if errors.Is(err, authentication.ErrExpiredClaim) ||
					errors.Is(err, authentication.ErrBlockedClaim) {
					logger.Info(
						"request blocked by middleware. path: %s, method: %s, middleware: auth. cause: %s",
						r.URL.Path, r.Method, err.Error())
				} else {
					logger.Error("jwt verification failed. jwt header: %s", err, jwtHeader)
				}
				http.Error(w, "verification error", http.StatusUnauthorized)
				return
			}
			if claim == nil {
				logger.Error("no claim found in jwt. jwt header: %s", err, jwtHeader)
				http.Error(w, "verification error", http.StatusUnauthorized)
				return
			}
			// set in context
			ctx := r.Context()
			ctx = context.WithValue(ctx, authentication.UserContextKey, claim.User.DeepCopy())
			ctx = context.WithValue(ctx, authentication.UserIDKey, claim.User.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRoleAccessMiddleware
//
//	@param allowdPermissionBit
//	@return mux.MiddlewareFunc
func GetRoleAccessMiddleware(allowdPermissionBit []uint64) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user context
			cval := r.Context().Value(authentication.UserContextKey)
			if cval == nil {
				err := errors.WithStack(ErrUserContextNotFound)
				logger.Error("failed to check role. path: %s, method: %s, middleware: role",
					err, r.URL.Path, r.Method)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			userContext, ok := cval.(authentication.UserContext)
			if !ok {
				err := errors.WithStack(ErrUserContextNotFound)
				logger.Error("failed to convert user context value", err)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			// verify
			ok, err := authorization.HasAuthority(userContext.PermissionBit, allowdPermissionBit)
			if err != nil {
				logger.Error("failed to intersept request. path: %s, method: %s, middleware: access", err, r.URL.Path, r.Method)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			if !ok {
				logger.Info("request blocked by middleware. path: %s, method: %s, middleware: access", r.URL.Path, r.Method)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// NotMatched
//
//	@param w
//	@param r
//	@return bool is handled
func NotMatched(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodOptions {
		logger.Info("handle unmatched request. path: %s, method: %s", r.URL.Path, r.Method)
		return false
	}
	logger.Info("handle unmatched cors request. path: %s, method: %s", r.URL.Path, r.Method)
	writeCors(w)
	return true
}

// writeCors
//
//	@param w
func writeCors(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
}
