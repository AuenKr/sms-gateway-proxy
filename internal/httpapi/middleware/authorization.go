package middleware

import (
	"net/http"

	"sms-gateway/internal/config"
	"sms-gateway/internal/httpapi"
)

func NewAuthorization(cfg config.Config) httpapi.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.HasServerAuthorization(r.Header.Get("Authorization")) {
				httpapi.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
