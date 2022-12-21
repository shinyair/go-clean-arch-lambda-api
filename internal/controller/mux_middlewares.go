package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"local.com/go-clean-lambda/internal/logger"
)

// GetLogMiddleware
//
//	@return mux.MiddlewareFunc
func GetLogMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("handle request. path: %s, method: %s, jwt: %s, middleware: log", r.URL.Path, r.Method, "TODO:")
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
