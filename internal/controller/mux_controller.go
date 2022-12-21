package controller

import (
	"net/http"
	"strings"

	"local.com/go-clean-lambda/internal/logger"

	"github.com/gorilla/mux"
)

// MuxRouterHandler
type MuxRouterHandler struct {
	mdfs  []mux.MiddlewareFunc
	hfunc func(w http.ResponseWriter, r *http.Request)
}

// ServeHTTP
//  @receiver h
//  @param w
//  @param r
func (h *MuxRouterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.hfunc(w, r)
}

// WithMiddlewares
//  @receiver h
//  @param mdfs
//  @return *MuxRouterHandler
func (h *MuxRouterHandler) WithMiddlewares(mdfs []mux.MiddlewareFunc) *MuxRouterHandler {
	h.mdfs = mdfs
	return h
}

// WithHandleFunc
//  @receiver h
//  @param hfunc
//  @return *MuxRouterHandler
func (h *MuxRouterHandler) WithHandleFunc(hfunc func(w http.ResponseWriter, r *http.Request)) *MuxRouterHandler {
	h.hfunc = hfunc
	return h
}

// GetMiddlewareFuncs
//  @receiver m
//  @return []mux.MiddlewareFunc
func (m *MuxRouterHandler) GetMiddlewareFuncs() []mux.MiddlewareFunc {
	return m.mdfs
}

// GetHandleFunc
//  @receiver h
//  @return w
//  @return r
//  @return func(w http.ResponseWriter, r *http.Request)
func (h *MuxRouterHandler) GetHandleFunc() func(w http.ResponseWriter, r *http.Request) {
	return h.hfunc
}

type MuxController Controller[*MuxRouterHandler]

// AddMuxRouter
//  @param c
//  @param path
//  @param methods
//  @param mdfs
//  @param handler
func AddMuxRouter(c MuxController, path string, methods []string, mdfs []mux.MiddlewareFunc, handler ErrorableHandler) {
	// check path
	path = formatPath(path)
	// get handler map
	handlerMap := c.GetHandlers()
	methodMap, ok := handlerMap[path]
	if !ok {
		methodMap = make(map[string]*MuxRouterHandler)
	}
	// build new filter chain
	h := &MuxRouterHandler{}
	h = h.WithHandleFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			logger.Error("failed to run request handler. path: %s, method: %s.", err, r.URL.Path, r.Method)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	})
	h = h.WithMiddlewares(mdfs)
	// set in map
	for _, method := range methods {
		methodMap[method] = h
	}
	handlerMap[path] = methodMap
}

func formatPath(path string) string {
	path = strings.TrimSpace(path)
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	return path
}
