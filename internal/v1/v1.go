package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kerraform/kerranamodb/internal/auth"
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
	auth   auth.Authenticator
	dmu    *dlock.DMutex
	driver driver.Driver
	logger *zap.Logger
	url    *url.URL
}

type HandlerConfig struct {
	Auth   auth.Authenticator
	Dmu    *dlock.DMutex
	Driver driver.Driver
	Logger *zap.Logger
	URL    string
}

func New(cfg *HandlerConfig) (*Handler, error) {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	return &Handler{
		auth:   cfg.Auth,
		dmu:    cfg.Dmu,
		driver: cfg.Driver,
		logger: cfg.Logger.Named("v1"),
		url:    u,
	}, nil
}

type CreateTenantRequest struct {
	Table string `json:"table"`
}

type CreateTenantResponse struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

func (h *Handler) CreateTenant() http.Handler {
	return handler.NewHandler(func(w http.ResponseWriter, r *http.Request) error {
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
			return nil
		case http.MethodPost:
		default:
			return kerrors.Wrap(errors.New("method not allowed"), kerrors.WithBadRequest("method not allowed"))
		}

		var req CreateTenantRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		tenant, err := h.driver.GetTenant(r.Context(), req.Table)
		if err != nil {
			if !errors.Is(err, driver.ErrTenantNotFound) {
				return err
			}
		}

		if tenant != nil {
			return kerrors.Wrap(err, kerrors.WithBadRequest(fmt.Sprintf("%s table already exists", req.Table)))
		}

		st, err := h.auth.Generate(r.Context(), &auth.Claims{
			Table: req.Table,
		})
		if err != nil {
			return err
		}

		if err := h.driver.CreateTenant(r.Context(), req.Table, st); err != nil {
			return err
		}

		q := h.url.Query()
		q.Set(auth.TokenQueryKey, st)
		h.url.RawQuery = q.Encode()

		return json.NewEncoder(w).Encode(&CreateTenantResponse{
			URL:   h.url.String(),
			Token: st,
		})
	})
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

	c, err := auth.FromContext(r.Context())
	if err != nil {
		return kerrors.Wrap(err, kerrors.WithInternalServerError())
	}

	if !c.TableAccessible(i.TableName) {
		return kerrors.Wrap(err, kerrors.WithForbidden("table not accessible"))
	}

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

	c, err := auth.FromContext(r.Context())
	if err != nil {
		return kerrors.Wrap(err, kerrors.WithInternalServerError())
	}

	if !c.TableAccessible(i.TableName) {
		return kerrors.Wrap(err, kerrors.WithForbidden("table not accessible"))
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

	c, err := auth.FromContext(r.Context())
	if err != nil {
		return kerrors.Wrap(err, kerrors.WithInternalServerError())
	}

	if !c.TableAccessible(i.TableName) {
		return kerrors.Wrap(err, kerrors.WithForbidden("table not accessible"))
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
