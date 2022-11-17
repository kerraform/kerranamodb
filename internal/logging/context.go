package logging

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

const (
	Key = "logger"
)

func FromCtx(ctx context.Context) (*zap.Logger, error) {
	l, ok := ctx.Value(Key).(*zap.Logger)
	if !ok {
		return nil, errors.New("failed to get logger from context")
	}
	return l, nil
}
