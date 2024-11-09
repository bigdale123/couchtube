package middleware

import (
	"encoding/json"
	"net/http"
	"slices"
)

func ReadOnlyGuard(next http.HandlerFunc) http.HandlerFunc {
	restrictedHttpMethods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(restrictedHttpMethods, r.Method) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"error": "Server is in read-only mode"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
