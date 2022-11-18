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
			var uid, t, k string
			switch v := req.Any().(type) {
			case *lockv1.LockRequest:
				uid = v.GetUid()
				t = v.GetTable()
				k = v.GetKey()
			case *lockv1.RLockRequest:
				uid = v.GetUid()
				t = v.GetTable()
				k = v.GetKey()
			case *lockv1.RUnlockRequest:
				uid = v.GetUid()
				t = v.GetTable()
				k = v.GetKey()
			case *lockv1.UnlockRequest:
				uid = v.GetUid()
				t = v.GetTable()
				k = v.GetKey()
			}

			l := logger.With(
				zap.String("uid", uid),
				zap.String("table", t),
				zap.String("key", k),
			)

			ctx = context.WithValue(ctx,
				logging.Key,
				l,
			)

			l.Info("access to lock node")
			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
