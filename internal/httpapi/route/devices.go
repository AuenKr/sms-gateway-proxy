package route

import (
	"net/http"

	"sms-gateway/internal/gateway"
	"sms-gateway/internal/httpapi"
)

const devicesPath = "/3rdparty/v1/devices"

type DevicesEndpoint struct {
	gateway     *gateway.Client
	middlewares []httpapi.Middleware
}

func NewDevicesEndpoint(gatewayClient *gateway.Client, authorization httpapi.Middleware) httpapi.Route {
	return &DevicesEndpoint{gateway: gatewayClient, middlewares: []httpapi.Middleware{authorization}}
}

func (e *DevicesEndpoint) Register(mux *http.ServeMux) {
	httpapi.Handle(mux, "GET /devices", http.HandlerFunc(e.listDevices), e.Middlewares())
}

func (e *DevicesEndpoint) Middlewares() []httpapi.Middleware {
	return e.middlewares
}

func (e *DevicesEndpoint) listDevices(w http.ResponseWriter, r *http.Request) {
	response, err := e.gateway.Do(r.Context(), http.MethodGet, devicesPath, r.URL.Query(), nil)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	httpapi.WriteGatewayResponse(w, response)
}
