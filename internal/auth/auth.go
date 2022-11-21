package auth

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kerraform/kerranamodb/internal/driver"
	"go.uber.org/zap"
)

const (
	loggerName = "authenticator"
)

type Authenticator interface {
	Generate(context.Context, *Claims) (string, error)
	Verify(context.Context, string) (*Claims, error)
}

type auth struct {
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey

	driver driver.Driver
	logger *zap.Logger
}

var _ Authenticator = (*auth)(nil)

func NewAuth(privateKeyPath, publicKeyPath string, d driver.Driver, logger *zap.Logger) (Authenticator, error) {
	prb, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	pub, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	prk, err := jwt.ParseEdPrivateKeyFromPEM(prb)
	if err != nil {
		return nil, err
	}

	puk, err := jwt.ParseEdPublicKeyFromPEM(pub)
	if err != nil {
		return nil, err
	}

	return &auth{
		driver:     d,
		logger:     logger.Named(loggerName),
		privateKey: prk,
		publicKey:  puk,
	}, nil
}

func (a *auth) Generate(ctx context.Context, claim *Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claim)
	st, err := t.SignedString(a.privateKey)
	if err != nil {
		return "", err
	}

	return st, err
}

func (a *auth) Verify(ctx context.Context, st string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(st, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.publicKey.(ed25519.PublicKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*Claims); ok && t.Valid {
		return claims, nil
	}

	return nil, errors.New("failed to verify token")
}
