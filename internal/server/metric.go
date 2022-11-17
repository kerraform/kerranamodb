package server

import "github.com/prometheus/client_golang/prometheus/promhttp"

func (s *Server) registerMetricsHandler() {
	s.mux.Handle("/metrics", promhttp.Handler())
}
