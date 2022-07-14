package controller

import (
	"net/http"

	"local.com/go-clean-lambda/internal/logger"

	"github.com/gorilla/mux"
)

// NewRouter
//  @param controllers
//  @return mux.Router
func NewRouter(controllers []MuxController) *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := NotMatched(w, r)
		if !ok {
			http.NotFound(w, r)
		}
	})
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := NotMatched(w, r)
		if !ok {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})
	for _, c := range controllers {
		s := r.PathPrefix(c.GetRootPath()).Subrouter()
		handlerMap := c.GetHandlers()
		for path, methodMap := range handlerMap {
			for method, handler := range methodMap {
				logger.Info("add request router. path: %s%s, method: %s", c.GetRootPath(), path, method)
				ss := s.Methods(method).Subrouter()
				ss.HandleFunc(path, handler.GetHandleFunc())
				for _, mdw := range handler.GetMiddlewareFuncs() {
					ss.Use(mdw)
				}
			}
		}
	}
	return r
}
