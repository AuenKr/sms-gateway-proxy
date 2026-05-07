package route

import (
	"net/http"

	"sms-gateway/internal/gateway"
	"sms-gateway/internal/httpapi"
)

type ProxyEndpoint struct {
	gateway     *gateway.Client
	middlewares []httpapi.Middleware
}

func NewProxyEndpoint(gatewayClient *gateway.Client, authorization httpapi.Middleware) httpapi.Route {
	return &ProxyEndpoint{gateway: gatewayClient, middlewares: []httpapi.Middleware{authorization}}
}

func (e *ProxyEndpoint) Register(mux *http.ServeMux) {
	httpapi.Handle(mux, "/", http.HandlerFunc(e.proxy), e.Middlewares())
}

func (e *ProxyEndpoint) Middlewares() []httpapi.Middleware {
	return e.middlewares
}

func (e *ProxyEndpoint) proxy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	response, err := e.gateway.Forward(r.Context(), r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httpapi.WriteGatewayResponse(w, response)
}
