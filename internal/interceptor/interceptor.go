package interceptor

import (
	"context"

	"github.com/bufbuild/connect-go"
	lockv1 "github.com/kerraform/kerranamodb/internal/gen/lock/v1"
	"github.com/kerraform/kerranamodb/internal/logging"
	"go.uber.org/zap"
)

func NewLoggingInterceptor(logger *zap.Logger) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			var method, uid, t, k string
			switch v := req.Any().(type) {
			case *lockv1.LockRequest:
				method = "lock"
				uid = v.GetUid()
				t = v.GetTable()
				k = v.GetKey()
			case *lockv1.RLockRequest:
				method = "rlock"
				uid = v.GetUid()
				t = v.GetTable()
				k = v.GetKey()
			case *lockv1.RUnlockRequest:
				method = "runlock"
				uid = v.GetUid()
				t = v.GetTable()
				k = v.GetKey()
			case *lockv1.UnlockRequest:
				method = "unlock"
				uid = v.GetUid()
				t = v.GetTable()
				k = v.GetKey()
			}

			l := logger.With(
				zap.String("method", method),
				zap.String("uid", uid),
				zap.String("table", t),
				zap.String("key", k),
			)

			ctx = context.WithValue(ctx,
				logging.Key,
				l,
			)

			l.Info("received request to lock node")
			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
