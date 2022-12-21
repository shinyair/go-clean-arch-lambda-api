package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"local.com/go-clean-lambda/internal/logger"
)

const (
	pong string = "pong"
)

// PingController
// works as extends MuxControllerImpl.
type PingController struct {
	*MuxControllerImpl
}

// NewPingController
//
//	@param logMdf
//	@param authMdf
//	@param roleMdf
//	@return *PingController
func NewPingController(
	logMdf mux.MiddlewareFunc,
	authMdf mux.MiddlewareFunc,
	roleMdf mux.MiddlewareFunc,
) *PingController {
	c := &PingController{
		MuxControllerImpl: NewMuxControllerImpl(
			"/api/ping",
			make(map[string]map[string]*MuxRouterHandler),
		),
	}
	c.AddMuxRouter("", []string{
		http.MethodGet,
	}, []mux.MiddlewareFunc{
		logMdf,
	}, func(w http.ResponseWriter, r *http.Request) error {
		// curl {host}/api/ping
		logger.Info("log middleware only")
		return c.WriteResponse(w, pong)
	})
	c.AddMuxRouter("", []string{
		http.MethodPost,
	}, []mux.MiddlewareFunc{
		logMdf,
		authMdf,
	}, func(w http.ResponseWriter, r *http.Request) error {
		// curl -X POST {host}/api/ping  {jwt}"
		logger.Info("log middleware & login middleware")
		return c.WriteResponse(w, pong)
	})
	c.AddMuxRouter("", []string{
		http.MethodPut,
	}, []mux.MiddlewareFunc{
		logMdf,
		authMdf,
		roleMdf,
	}, func(w http.ResponseWriter, r *http.Request) error {
		// curl -X PUT {host}/api/ping -H "Authorization: Bearer {jwt}"
		logger.Info("log middleware & login middleware & role middleware")
		return c.WriteResponse(w, pong)
	})
	return c
}
