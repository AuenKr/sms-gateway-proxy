package middleware

import (
	"crypto/subtle"
	"net/http"

	"sms-gateway/internal/config"
	"sms-gateway/internal/httpapi"
)

func NewAuthorization(cfg config.Config) httpapi.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := r.Header.Get("Authorization")
			if subtle.ConstantTimeCompare([]byte(value), []byte(cfg.ServerAuthKey)) != 1 {
				httpapi.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
