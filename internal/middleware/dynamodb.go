package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type amazonAPIContextKey string

const (
	amzTargetKey = "X-Amz-Target"

	AmazonAPIVersionKey   amazonAPIContextKey = "AmazonAPIVersionKey"
	AmazonAPIOperationKey amazonAPIContextKey = "AmazonAPIOperationKey"
)

func DynamoDB() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := r.Header.Get(amzTargetKey)
			if v == "" {
				err := fmt.Errorf("failed get `X-Amz-Target` HTTP header")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			el := strings.Split(v, ".")
			if len(el) != 2 {
				err := fmt.Errorf("malformed `X-Amz-Target`")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, AmazonAPIVersionKey, el[0])
			ctx = context.WithValue(ctx, AmazonAPIOperationKey, el[1])

			req := r.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	}
}
