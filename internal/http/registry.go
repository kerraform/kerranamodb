package http

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

	v1 := s.mux.PathPrefix(v1Path).Subrouter()
	v1.Use(middleware.DynamoDB())

	// ProvisionTenants
	v1.Methods(http.MethodPost).Path("/tenant").Handler(s.v1.CreateTenant())

	if v := s.auth; v != nil {
		v1.Use(middleware.Auth(v))
	}

	// Note(KeisukeYamashita):
	// Paths can be configured by `dynamodb_endpoint` field on developer side.
	// Thus, for future development, I will version-ize this API server.
	v1.Methods(http.MethodPost).Path("/").Handler(s.v1.Handler())
}
