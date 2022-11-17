package interceptor

import (
	"context"

	"github.com/bufbuild/connect-go"
	"go.uber.org/zap"
)

func NewLoggingInterceptor(logger *zap.Logger) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			logger.Info("access to grpc")
			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
