package http

import (
	"errors"
	"net/http"

	kerrors "github.com/kerraform/kerranamodb/internal/errors"
	"github.com/kerraform/kerranamodb/internal/handler"
	"github.com/kerraform/kerranamodb/internal/middleware"
)

func (s *Server) registerUtilHandler() {
	s.mux.Methods(http.MethodGet).Path("/healthz").Handler(s.HealthCheck())
	s.mux.NotFoundHandler = middleware.AccessLog(s.logger)(s.NotFound())
}

func (s *Server) HealthCheck() http.Handler {
	return handler.NewHandler(func(w http.ResponseWriter, _ *http.Request) error {
		if s.dmu.Ready {
			w.WriteHeader(http.StatusOK)
			return nil
		}

		return kerrors.Wrap(errors.New("dlock not ready"), kerrors.WithInternalServerError())
	})
}

func (s *Server) NotFound() http.Handler {
	return handler.NewHandler(func(w http.ResponseWriter, _ *http.Request) error {
		w.WriteHeader(http.StatusNotFound)
		return errors.New("not found")
	})
}
