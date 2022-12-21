package controller

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"local.com/go-clean-lambda/internal/logger"
)

// MuxRouterHandler implements interface http.Handler.
type MuxRouterHandler struct {
	mdfs  []mux.MiddlewareFunc
	hfunc func(w http.ResponseWriter, r *http.Request)
}

// ServeHTTP
//
//	@receiver h
//	@param w
//	@param r
func (h *MuxRouterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.hfunc(w, r)
}

// WithMiddlewares
//
//	@receiver h
//	@param mdfs
//	@return *MuxRouterHandler
func (h *MuxRouterHandler) WithMiddlewares(mdfs []mux.MiddlewareFunc) *MuxRouterHandler {
	h.mdfs = mdfs
	return h
}

// WithHandleFunc
//
//	@receiver h
//	@param hfunc
//	@return *MuxRouterHandler
func (h *MuxRouterHandler) WithHandleFunc(hfunc func(w http.ResponseWriter, r *http.Request)) *MuxRouterHandler {
	h.hfunc = hfunc
	return h
}

// GetMiddlewareFuncs
//
//	@receiver m
//	@return []mux.MiddlewareFunc
func (h *MuxRouterHandler) GetMiddlewareFuncs() []mux.MiddlewareFunc {
	return h.mdfs
}

// GetHandleFunc
//
//	@receiver h
//	@return w
//	@return r
//	@return func(w http.ResponseWriter, r *http.Request)
func (h *MuxRouterHandler) GetHandleFunc() func(w http.ResponseWriter, r *http.Request) {
	return h.hfunc
}

// MuxController extends interface Controller[T http.Handler].
type MuxController Controller[*MuxRouterHandler]

// MuxControllerImpl implements interface MuxController.
type MuxControllerImpl struct {
	rootPath   string
	handlerMap map[string]map[string]*MuxRouterHandler
}

// NewMuxControllerImpl
//
//	@param rootPath
//	@param handlerMap
//	@return *MuxControllerImpl
func NewMuxControllerImpl(
	rootPath string,
	handlerMap map[string]map[string]*MuxRouterHandler,
) *MuxControllerImpl {
	return &MuxControllerImpl{
		rootPath:   rootPath,
		handlerMap: handlerMap,
	}
}

// GetRootPath
//
//	@receiver c
//	@return string
func (c *MuxControllerImpl) GetRootPath() string {
	return c.rootPath
}

// GetHandlers
//
//	@receiver c
//	@return map
func (c *MuxControllerImpl) GetHandlers() map[string]map[string]*MuxRouterHandler {
	return c.handlerMap
}

// AddMuxRouter
//
//	@param c
//	@param path
//	@param methods
//	@param mdfs
//	@param handler
func (c *MuxControllerImpl) AddMuxRouter(
	path string,
	methods []string,
	mdfs []mux.MiddlewareFunc,
	handler ErrorableHandler,
) {
	// check path
	path = c.formatPath(path)
	// get handler map
	handlerMap := c.GetHandlers()
	methodMap, ok := handlerMap[path]
	// init map
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

// formatPath
//
// make sure the path begins with '/' and not ends with '/'
//
//	@receiver c
//	@param path
//	@return string
func (c *MuxControllerImpl) formatPath(path string) string {
	path = strings.TrimSpace(path)
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	return path
}

// WriteResponse
//
//	@receiver c
//	@param w
//	@param body
//	@return error
func (c *MuxControllerImpl) WriteResponse(w http.ResponseWriter, body string) error {
	_, err := w.Write([]byte(body))
	if err != nil {
		return errors.Wrap(errors.New(err.Error()), "write response failed")
	}
	return nil
}
