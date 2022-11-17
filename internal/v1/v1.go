package v1

import (
	"net/http"

	"github.com/kerraform/kerranamodb/internal/driver"
	"github.com/kerraform/kerranamodb/internal/handler"
	"go.uber.org/zap"
)

type DataType string

const (
	DataTypeAddGPGKey DataType = "gpg-keys"
)

type Handler struct {
	logger *zap.Logger
	driver driver.Driver
}

type HandlerConfig struct {
	Driver driver.Driver
	Logger *zap.Logger
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		driver: cfg.Driver,
		logger: cfg.Logger.Named("v1"),
	}
}

func (h *Handler) Handler() http.Handler {
	return handler.NewHandler(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})
}
