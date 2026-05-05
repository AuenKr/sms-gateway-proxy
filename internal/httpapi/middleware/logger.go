package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"sms-gateway/internal/httpapi"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "sms-gateway ", log.LstdFlags|log.LUTC)
}

func NewRequestLogging(logger *log.Logger) httpapi.OrderedMiddleware {
	return httpapi.OrderedMiddleware{
		Order: 1,
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				startedAt := time.Now()
				recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
				next.ServeHTTP(recorder, r)
				logger.Printf("method=%s path=%s status=%d remote_addr=%s duration=%s", r.Method, r.URL.Path, recorder.statusCode, r.RemoteAddr, time.Since(startedAt).String())
			})
		},
	}
}
