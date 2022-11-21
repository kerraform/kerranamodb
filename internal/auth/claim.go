package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

const (
	contextKey = "claims"
)

type Claims struct {
	*jwt.RegisteredClaims
	Table string `json:"table"`
}

func FromContext(ctx context.Context) (*Claims, error) {
	c, ok := ctx.Value(contextKey).(*Claims)
	if !ok {
		return nil, errors.New("failed to get claims from context")
	}

	return c, nil
}

func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, contextKey, claims)
}

func (c *Claims) TableAccessible(table string) bool {
	return c.Table == table
}

func (c *Claims) Valid() error {
	return nil
}
