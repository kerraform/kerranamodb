package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kerraform/kerranamodb/internal/errors"
	"github.com/kerraform/kerranamodb/internal/logging"
	"go.uber.org/zap"
)

type Error struct {
	Message string `json:"message"`
}

type Handler struct {
	HandleFunc HandlerFunc
}

// HandlerFunc represents the registry handler
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP Implements the http.Handler
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hErr := h.HandleFunc(w, r)
	if hErr == nil {
		return
	}

	l, err := logging.FromCtx(r.Context())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get logger; %v", err)
		return
	}

	l.Error("error to response", zap.Error(hErr))
	if err := errors.ServeJSON(w, hErr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewHandler(f HandlerFunc) http.Handler {
	return &Handler{
		HandleFunc: f,
	}
}
