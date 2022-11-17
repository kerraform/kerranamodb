package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kerraform/kerranamodb/internal/driver"
	"github.com/kerraform/kerranamodb/internal/dynamodb"
	"github.com/kerraform/kerranamodb/internal/dynamodb/api"
	"github.com/kerraform/kerranamodb/internal/errors"
	"github.com/kerraform/kerranamodb/internal/handler"
	"github.com/kerraform/kerranamodb/internal/middleware"
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
		method := dynamodb.OperationType(r.Context().Value(middleware.AmazonAPIOperationKey).(string))

		switch method {
		case dynamodb.OperationTypeDeleteItem:
			var i api.DeleteInput

			if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
				return err
			}
			defer r.Body.Close()

			lid, err := i.GetLockID()
			if err != nil {
				return errors.Wrap(err, errors.WithBadRequest("failed to get lock id"))
			}

			return h.driver.DeleteLock(r.Context(), i.TableName, lid)
		case dynamodb.OperationTypeGetItem:
			var i api.UpsertInput

			if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
				return err
			}
			defer r.Body.Close()

			lid, err := i.GetLockID()
			if err != nil {
				return errors.Wrap(err, errors.WithBadRequest("failed to get lock id"))
			}

			info, err := h.driver.GetLock(r.Context(), i.TableName, lid)
			if err != nil {
				return err
			}

			r := &api.UpsertInput{
				TableName: i.TableName,
				Item: map[string]map[string]string{
					api.InfoKey: {
						api.SKey: string(info),
					},
					api.LockIDKey: {
						api.SKey: string(lid),
					},
				},
			}

			return json.NewEncoder(w).Encode(r)
		case dynamodb.OperationTypePutItem:
			var i api.UpsertInput

			if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
				return err
			}
			defer r.Body.Close()

			info, err := i.GetInfo()
			if err != nil {
				return errors.Wrap(err, errors.WithBadRequest("failed to get info"))
			}

			if info == "" {
				return errors.Wrap(err, errors.WithBadRequest("empty info"))
			}

			lid, err := i.GetLockID()
			if err != nil {
				return errors.Wrap(err, errors.WithBadRequest("failed to get lock id"))
			}

			if err := h.driver.SaveLock(r.Context(), i.TableName, lid, driver.Info(info)); err != nil {
				return err
			}

			return nil
		default:
			err := fmt.Errorf("method: %s not allowed", method)
			return errors.Wrap(err, errors.WithBadRequest(err.Error()))
		}
	})
}
