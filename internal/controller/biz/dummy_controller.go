package bizcontroller

import (
	"fmt"
	"net/http"

	"local.com/go-clean-lambda/internal/controller"
	"local.com/go-clean-lambda/internal/logger"
	"local.com/go-clean-lambda/internal/usecase"

	"github.com/gorilla/mux"
)

// DummyController
type DummyController struct {
	rootPath   string
	handlerMap map[string]map[string]*controller.MuxRouterHandler
	usecase    *usecase.DummyUseCase
}

// GetRootPath
//  @receiver c
//  @return string
func (c *DummyController) GetRootPath() string {
	return c.rootPath
}

// GetHandlers
//  @receiver c
//  @return map
func (c *DummyController) GetHandlers() map[string]map[string]*controller.MuxRouterHandler {
	return c.handlerMap
}

// NewDummyController
//  @param logMdf
//  @param usecase
//  @return controller.MuxController
func NewDummyController(logMdf mux.MiddlewareFunc, usecase *usecase.DummyUseCase) controller.MuxController {
	c := &DummyController{
		rootPath:   "/api/dummy",
		handlerMap: make(map[string]map[string]*controller.MuxRouterHandler),
		usecase:    usecase,
	}
	controller.AddMuxRouter(c, "/{id}", []string{
		http.MethodGet,
	}, []mux.MiddlewareFunc{
		logMdf,
	}, func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		return c.handleGet(w, r, vars["id"])
	})
	controller.AddMuxRouter(c, "", []string{
		http.MethodPost,
	}, []mux.MiddlewareFunc{
		logMdf,
	}, func(w http.ResponseWriter, r *http.Request) error {
		return c.handlePost(w, r)
	})
	controller.AddMuxRouter(c, "/{id}", []string{
		http.MethodDelete,
	}, []mux.MiddlewareFunc{
		logMdf,
	}, func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		return c.handleDelete(w, r, vars["id"])
	})
	return c
}

// handleGet
//  @receiver c
//  @param w
//  @param r
//  @param id
//  @return error
func (c *DummyController) handleGet(w http.ResponseWriter, r *http.Request, id string) error {
	logger.Debug("get by id: %s", id)
	bo, err := c.usecase.Get(r.Context(), id)
	if err != nil {
		return err
	}
	s := fmt.Sprintf("handle get. id: %s, got bo: %s", id, logger.Pretty(bo))
	logger.Info(s)
	w.Write([]byte(s))
	return nil
}

// handlePost
//  @receiver c
//  @param w
//  @param r
//  @return error
func (c *DummyController) handlePost(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	attr := r.URL.Query().Get("attr")
	logger.Debug("handle post. id: %s, name: %s", id, name)
	bo, err := c.usecase.Add(r.Context(), &usecase.DummyBo{
		ID:   id,
		Name: name,
		Attr: attr,
	})
	if err != nil {
		return err
	}
	s := logger.Pretty(bo)
	logger.Debug("handle add. bo: %s", s)
	w.Write([]byte(s))
	return nil
}

// handleDelete
//  @receiver c
//  @param w
//  @param r
//  @param id
//  @return error
func (c *DummyController) handleDelete(w http.ResponseWriter, r *http.Request, id string) error {
	logger.Debug("delete by id: %s", id)
	bo, err := c.usecase.Remove(r.Context(), id)
	if err != nil {
		return err
	}
	if bo == nil {
		return fmt.Errorf("failed to handle delete. id: %s", id)
	}
	s := fmt.Sprintf("handle delete. id: %s", id)
	logger.Debug(s)
	w.Write([]byte(s))
	return nil
}
