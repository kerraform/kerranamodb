package middleware

import (
	"fmt"
	"net/http"
)

func JSON() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("hohoh")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			next.ServeHTTP(w, r)
		})
	}
}
