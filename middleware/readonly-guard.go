package middleware

import (
	"net/http"
)

func ReadOnlyGuard(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Server is in read-only mode"))

		return
	})
}
