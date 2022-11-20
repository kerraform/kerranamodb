package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kerraform/kerranamodb/internal/dlock"
	"github.com/kerraform/kerranamodb/internal/driver"
	"github.com/kerraform/kerranamodb/internal/dynamodb"
	"github.com/kerraform/kerranamodb/internal/dynamodb/api"
	kerrors "github.com/kerraform/kerranamodb/internal/errors"
	"github.com/kerraform/kerranamodb/internal/handler"
	"github.com/kerraform/kerranamodb/internal/middleware"
	"go.uber.org/zap"
)

type Handler struct {
	dmu    *dlock.DMutex
	logger *zap.Logger
	driver driver.Driver
}

type HandlerConfig struct {
	Dmu    *dlock.DMutex
	Driver driver.Driver
	Logger *zap.Logger
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		dmu:    cfg.Dmu,
		driver: cfg.Driver,
		logger: cfg.Logger.Named("v1"),
	}
}

func (h *Handler) Handler() http.Handler {
	return handler.NewHandler(func(w http.ResponseWriter, r *http.Request) error {
		method := dynamodb.OperationType(r.Context().Value(middleware.AmazonAPIOperationKey).(string))

		switch method {
		case dynamodb.OperationTypeDeleteItem:
			return h.deleteLock(w, r)
		case dynamodb.OperationTypeGetItem:
			return h.getLock(w, r)
		case dynamodb.OperationTypePutItem:
			return h.putLock(w, r)
		default:
			err := fmt.Errorf("method: %s not allowed", method)
			return kerrors.Wrap(err, kerrors.WithBadRequest(err.Error()))
		}
	})
}

func (h *Handler) deleteLock(_ http.ResponseWriter, r *http.Request) error {
	var i api.DeleteInput

	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		return err
	}
	defer r.Body.Close()

	lid, err := i.GetLockID()
	if err != nil {
		return kerrors.Wrap(err, kerrors.WithBadRequest("failed to get lock id"))
	}

	return h.driver.DeleteLock(r.Context(), i.TableName, lid)
}

func (h *Handler) getLock(w http.ResponseWriter, r *http.Request) error {
	var i api.GetInput

	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		return err
	}
	defer r.Body.Close()

	lid, err := i.GetLockID()
	if err != nil {
		return kerrors.Wrap(err, kerrors.WithBadRequest("failed to get lock id"))
	}

	dlid := dlock.From(i.TableName, string(lid))

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := h.dmu.RLock(ctx, dlid); err != nil {
		h.logger.Error("someone in the cluster has the lock or trying to get it", zap.Error(err))
		return kerrors.Wrap(fmt.Errorf("state is locked"), kerrors.WithConditionalCheckFailedException())
	}
	defer h.dmu.RUnlock(ctx, dlid)

	info, err := h.driver.GetLock(r.Context(), i.TableName, lid)
	if err != nil {
		return err
	}

	res := &api.PutInput{
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

	return json.NewEncoder(w).Encode(res)
}

func (h *Handler) putLock(w http.ResponseWriter, r *http.Request) error {
	var i api.PutInput

	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		return err
	}
	defer r.Body.Close()

	info, err := i.GetInfo()
	if err != nil {
		return kerrors.Wrap(err, kerrors.WithBadRequest("failed to get info"))
	}

	if info == "" {
		return kerrors.Wrap(err, kerrors.WithBadRequest("empty info"))
	}

	lid, err := i.GetLockID()
	if err != nil {
		return kerrors.Wrap(err, kerrors.WithBadRequest("failed to get lock id"))
	}

	dlid := dlock.From(i.TableName, string(lid))

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.dmu.RLock(ctx, dlid); err != nil {
		h.logger.Error("someone in the cluster has the lock or trying to get it", zap.Error(err))
		return kerrors.Wrap(fmt.Errorf("state is locked"), kerrors.WithConditionalCheckFailedException())
	}

	hasLock, err := h.driver.HasLock(r.Context(), i.TableName, lid)
	if err != nil {
		return err
	}
	if hasLock {
		return kerrors.Wrap(fmt.Errorf("state is locked"), kerrors.WithConditionalCheckFailedException())
	}
	h.dmu.RUnlock(ctx, dlid)

	if err := h.dmu.Lock(r.Context(), dlid); err != nil {
		h.logger.Error("someone in the cluster has the lock or trying to get it", zap.Error(err))
		return kerrors.Wrap(fmt.Errorf("state is locked"), kerrors.WithConditionalCheckFailedException())
	}
	defer h.dmu.Unlock(ctx, dlid)

	if err := h.driver.SaveLock(r.Context(), i.TableName, lid, driver.Info(info)); err != nil {
		return err
	}

	return nil
}
