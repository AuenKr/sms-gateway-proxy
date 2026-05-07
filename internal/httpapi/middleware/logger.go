package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/fx"

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

type NewLoggerResult struct {
	fx.Out

	Logger *log.Logger
}

func NewLogger() NewLoggerResult {
	return NewLoggerResult{Logger: log.New(os.Stdout, "sms-gateway ", log.LstdFlags|log.LUTC)}
}

type NewRequestLoggingParams struct {
	fx.In

	Logger *log.Logger
}

type NewRequestLoggingResult struct {
	fx.Out

	Middleware httpapi.OrderedMiddleware `group:"global_middleware"`
}

func NewRequestLogging(in NewRequestLoggingParams) NewRequestLoggingResult {
	return NewRequestLoggingResult{
		Middleware: httpapi.OrderedMiddleware{
			Order: 1,
			Middleware: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					startedAt := time.Now()
					recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
					next.ServeHTTP(recorder, r)
					in.Logger.Printf("method=%s path=%s status=%d remote_addr=%s duration=%s", r.Method, r.URL.Path, recorder.statusCode, r.RemoteAddr, time.Since(startedAt).String())
				})
			},
		},
	}
}
