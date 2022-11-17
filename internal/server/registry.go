package server

import (
	"net/http"

	"github.com/kerraform/kerranamodb/internal/middleware"
)

const (
	v1Path = "/v1"
)

func (s *Server) registerRegistryHandler() {
	s.mux.Use(middleware.JSON())
	s.mux.Use(middleware.AccessLog(s.logger))
	s.mux.Use(middleware.AccessMetric(s.metric))
	s.mux.Use(middleware.DynamoDB())

	v1 := s.mux.PathPrefix(v1Path).Subrouter()

	// Note(KeisukeYamashita):
	// Paths can be configured by `dynamodb_endpoint` field.
	// Thus, for future development, I will version-ize this API server.
	v1.Methods(http.MethodPost).Path("").Handler(s.v1.Handler())
}
