package route

import (
	"net/http"

	"sms-gateway/internal/gateway"
	"sms-gateway/internal/httpapi"

	"go.uber.org/fx"
)

type ProxyEndpoint struct {
	gateway     *gateway.Client
	middlewares []httpapi.Middleware
}

type NewProxyEndpointParams struct {
	fx.In

	GatewayClient *gateway.Client
	Authorization httpapi.Middleware
}

type NewProxyEndpointResult struct {
	fx.Out

	Route httpapi.Route `group:"routes"`
}

func NewProxyEndpoint(in NewProxyEndpointParams) NewProxyEndpointResult {
	return NewProxyEndpointResult{
		Route: &ProxyEndpoint{
			gateway:     in.GatewayClient,
			middlewares: []httpapi.Middleware{in.Authorization},
		},
	}
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
