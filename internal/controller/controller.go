package controller

import (
	"net/http"
)

type Controller[T http.Handler] interface {
	GetRootPath() string
	GetHandlers() map[string]map[string]T
}

type ErrorableHandler func(w http.ResponseWriter, r *http.Request) error
