package dlock

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	lockv1 "github.com/kerraform/kerranamodb/internal/gen/lock/v1"
	"github.com/kerraform/kerranamodb/internal/gen/lock/v1/lockv1connect"
	"github.com/kerraform/kerranamodb/internal/interceptor"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	port int
	mux  *http.ServeMux
}

type LockService struct {
}

type LockServiceOptions struct {
	Port   int
	Logger *zap.Logger
}

func NewLockService(opts *LockServiceOptions) *Server {
	svc := &LockService{}

	interceptors := connect.WithInterceptors(
		interceptor.NewLoggingInterceptor(opts.Logger.Named("dlock")),
	)
	path, handler := lockv1connect.NewLockServiceHandler(svc, interceptors)
	s := &Server{
		port: opts.Port,
		mux:  http.NewServeMux(),
	}

	s.mux.Handle(path, handler)
	return s
}

func (s *Server) Serve() error {
	return http.ListenAndServe(
		fmt.Sprintf("localhost:%d", s.port),
		h2c.NewHandler(s.mux, &http2.Server{}),
	)
}

func (s *LockService) Lock(ctx context.Context, req *connect.Request[lockv1.LockRequest]) (*connect.Response[lockv1.LockResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&lockv1.LockResponse{})
	return res, nil
}

func (s *LockService) Unlock(ctx context.Context, req *connect.Request[lockv1.UnlockRequest]) (*connect.Response[lockv1.UnlockResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&lockv1.UnlockResponse{})
	return res, nil
}

func (s *LockService) RLock(ctx context.Context, req *connect.Request[lockv1.RLockRequest]) (*connect.Response[lockv1.RLockResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&lockv1.RLockResponse{})
	return res, nil
}

func (s *LockService) RUnlock(ctx context.Context, req *connect.Request[lockv1.RUnlockRequest]) (*connect.Response[lockv1.RUnlockResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&lockv1.RUnlockResponse{})
	return res, nil
}
