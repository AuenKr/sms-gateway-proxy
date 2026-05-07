package middleware

import (
	"crypto/subtle"
	"net/http"

	"sms-gateway/internal/config"
	"sms-gateway/internal/httpapi"

	"go.uber.org/fx"
)

type NewAuthorizationParams struct {
	fx.In

	Config config.Config
}

type NewAuthorizationResult struct {
	fx.Out

	Middleware httpapi.Middleware
}

func NewAuthorization(in NewAuthorizationParams) NewAuthorizationResult {
	return NewAuthorizationResult{
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				value := r.Header.Get("Authorization")
				if subtle.ConstantTimeCompare([]byte(value), []byte(in.Config.ServerAuthKey)) != 1 {
					httpapi.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
					return
				}
				next.ServeHTTP(w, r)
			})
		},
	}
}
