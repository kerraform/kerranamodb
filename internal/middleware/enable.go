package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kerraform/kerranamodb/internal/handler"
)

type RegistryType string

const (
	ModuleRegistryType   RegistryType = "module"
	ProviderRegistryType RegistryType = "provider"
)

func Enable(rType RegistryType, enabled bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if enabled {
				next.ServeHTTP(w, r)
				return
			}

			e := &handler.Error{
				Message: fmt.Sprintf("%s not enabled", rType),
			}

			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(e); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})
	}
}
