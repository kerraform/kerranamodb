package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kerraform/kerranamodb/internal/auth"
)

func Auth(a auth.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: This should be improved.
			if !strings.HasPrefix(r.URL.String(), "/v1/tables") {
				next.ServeHTTP(w, r)
				return
			}

			st := mux.Vars(r)["token"]
			if st == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			c, err := a.Verify(r.Context(), st)
			if err != nil {
				http.Error(w, "failed to authenticate", http.StatusUnauthorized)
				return
			}

			req := r.WithContext(auth.WithClaims(r.Context(), c))
			next.ServeHTTP(w, req)
		})
	}
}
