package middleware

import (
	"net/http"

	"sms-gateway/internal/httpapi"
)

const maxRequestBodyBytes int64 = 2 << 20 // 2 MiB

func NewContentLengthValidation() httpapi.OrderedMiddleware {
	return httpapi.OrderedMiddleware{
		Order: 2,
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if len(r.TransferEncoding) > 0 {
					httpapi.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "transfer encoding is not supported"})
					return
				}

				if r.ContentLength > maxRequestBodyBytes {
					httpapi.WriteJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "request body exceeds 2 MiB limit"})
					return
				}

				next.ServeHTTP(w, r)
			})
		},
	}
}
