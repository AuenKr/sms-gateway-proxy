package route

import (
	"net/http"
	"strings"

	"sms-gateway/internal/gateway"
	"sms-gateway/internal/httpapi"
)

const webhooksPath = "/3rdparty/v1/webhooks"

type WebhooksEndpoint struct {
	gateway     *gateway.Client
	middlewares []httpapi.Middleware
}

func NewWebhooksEndpoint(gatewayClient *gateway.Client, authorization httpapi.Middleware) httpapi.Route {
	return &WebhooksEndpoint{gateway: gatewayClient, middlewares: []httpapi.Middleware{authorization}}
}

func (e *WebhooksEndpoint) Register(mux *http.ServeMux) {
	httpapi.Handle(mux, "POST /webhooks", http.HandlerFunc(e.registerWebhook), e.Middlewares())
	httpapi.Handle(mux, "DELETE /webhooks/{id}", http.HandlerFunc(e.removeWebhook), e.Middlewares())
}

func (e *WebhooksEndpoint) Middlewares() []httpapi.Middleware {
	return e.middlewares
}

func (e *WebhooksEndpoint) registerWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req httpapi.RegisterWebhookRequest
	if err := httpapi.DecodeJSONStrict(r, &req); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	if strings.TrimSpace(req.URL) == "" {
		httpapi.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "url is required"})
		return
	}
	if strings.TrimSpace(req.Event) == "" {
		httpapi.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "event is required"})
		return
	}

	response, err := e.gateway.DoJSON(r.Context(), http.MethodPost, webhooksPath, r.URL.Query(), req)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	httpapi.WriteGatewayResponse(w, response)
}

func (e *WebhooksEndpoint) removeWebhook(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		httpapi.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "webhook id is required"})
		return
	}

	response, err := e.gateway.Do(r.Context(), http.MethodDelete, webhooksPath+"/"+id, r.URL.Query(), nil)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	httpapi.WriteGatewayResponse(w, response)
}
