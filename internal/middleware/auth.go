package middleware

import (
	"fmt"
	"net/http"

	"github.com/kerraform/kerranamodb/internal/auth"
)

func Auth(auth auth.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("hoho")
			next.ServeHTTP(w, r)
		})
	}
}
