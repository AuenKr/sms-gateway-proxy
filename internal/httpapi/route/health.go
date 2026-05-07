package route

import (
	"net/http"

	"sms-gateway/internal/httpapi"
)

type HealthEndpoint struct {
	middlewares []httpapi.Middleware
}

func NewHealthEndpoint() httpapi.Route {
	return &HealthEndpoint{}
}

func (e *HealthEndpoint) Register(mux *http.ServeMux) {
	httpapi.Handle(mux, "GET /health", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		httpapi.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}), e.Middlewares())
}

func (e *HealthEndpoint) Middlewares() []httpapi.Middleware {
	return e.middlewares
}
