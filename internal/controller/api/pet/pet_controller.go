package pet

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"local.com/go-clean-lambda/internal/controller"
	"local.com/go-clean-lambda/internal/logger"
)

// Pet example.
type Pet struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PetController
// works as extends MuxControllerImpl.
type APIPetController struct {
	*controller.MuxControllerImpl
}

// do not add param or return comments for functions that are not swagger operations
// it will cause parse error.
func NewPetController(logMdf mux.MiddlewareFunc) *APIPetController {
	c := &APIPetController{
		MuxControllerImpl: controller.NewMuxControllerImpl(
			"/v1/pets",
			make(map[string]map[string]*controller.MuxRouterHandler),
		),
	}
	c.AddMuxRouter("/{id}", []string{
		http.MethodGet,
	}, []mux.MiddlewareFunc{
		logMdf,
	}, func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		return c.handleGet(w, r, vars["id"])
	})
	return c
}

// handleGet godoc
// @Summary      get a pet
// @Description  get a pet by ID
// @Tags         pets
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Pet ID"
// @Param        name query     string  true  "Pet Name"
// @Success      200  {object}  Pet
// @Failure      400  {object}  string
// @Failure      404  {string}  string
// @Failure      500  {string}  string
// @Router       /v1/pets/{id} [get]
//
// get pet by id in path param and receive name as query param.
func (c *APIPetController) handleGet(w http.ResponseWriter, r *http.Request, id string) error {
	logger.Debug("get by id: %s", id)
	name := r.URL.Query().Get("name")
	pet := &Pet{
		ID:          id,
		Name:        name,
		Description: fmt.Sprintf("get by id: %s, query param name: %s, done", id, name),
	}
	_, err := w.Write([]byte(logger.Pretty(pet)))
	if err != nil {
		return errors.Wrap(errors.New(err.Error()), "write response failed")
	}
	return nil
}
