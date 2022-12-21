package car

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"local.com/go-clean-lambda/internal/controller"
	"local.com/go-clean-lambda/internal/logger"
)

type Car struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type APICarController struct {
	*controller.MuxControllerImpl
}

func NewCarController(logMdf mux.MiddlewareFunc) *APICarController {
	c := &APICarController{
		MuxControllerImpl: controller.NewMuxControllerImpl(
			"/v1/cars",
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
// @Summary      get a car
// @Description  get a car by ID
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Car ID, use 'error_id' to get a error response"
// @Param        type query     string  true  "Car Type"
// @Success      200  {object}  Car
// @Failure      400  {string}  string
// @Failure      404  {string}  string
// @Failure      500  {object}  controller.ErrorResponse
// @Router       /v1/cars/{id} [get]
//
// get pet by id in path param and receive name as query param.
func (c *APICarController) handleGet(w http.ResponseWriter, r *http.Request, id string) error {
	logger.Debug("get by id: %s", id)
	cartype := r.URL.Query().Get("type")
	var car *Car
	var err error
	if id != "error_id" {
		car = &Car{
			ID:          id,
			Type:        cartype,
			Description: fmt.Sprintf("get by id: %s, query param type: %s, done", id, cartype),
		}
	} else {
		err = errors.Errorf("mock get error. get by id: %s, query param type: %s", id, cartype)
	}
	var resp string
	if err != nil {
		resp = logger.Pretty(&controller.ErrorResponse{
			ErrorType:    "use case error",
			ErrorMessage: fmt.Sprintf("failed to get by id: %s", id),
		})
		// return errors.Wrap(err, "write response failed")
	} else {
		resp = logger.Pretty(car)
	}
	_, err = w.Write([]byte(resp))
	if err != nil {
		return errors.Wrap(errors.New(err.Error()), "write response failed")
	}
	return nil
}
