package auth

import (
	"context"
	"crypto"
	"io/ioutil"

	"github.com/golang-jwt/jwt"
	"github.com/kerraform/kerranamodb/internal/driver"
	"go.uber.org/zap"
)

const (
	loggerName = "authenticator"
)

type Authenticator interface {
	Generate(context.Context, *Claims) (string, error)
	Verify(context.Context, string)
}

type auth struct {
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey

	driver driver.Driver
	logger *zap.Logger
}

var _ Authenticator = (*auth)(nil)

func NewAuth(keypath string, d driver.Driver, logger *zap.Logger) (Authenticator, error) {
	b, err := ioutil.ReadFile(keypath)
	if err != nil {
		return nil, err
	}

	prk, err := jwt.ParseEdPrivateKeyFromPEM(b)
	if err != nil {
		return nil, err
	}

	puk, err := jwt.ParseEdPrivateKeyFromPEM(b)
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

func (a *auth) Verify(ctx context.Context, token string) {}
