package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"

	"sms-gateway/internal/config"
	"sms-gateway/internal/encryption"
	"sms-gateway/internal/gateway"
	"sms-gateway/internal/httpapi"
)

const messagesPath = "/messages"

type MessagesEndpoint struct {
	gateway     *gateway.Client
	encryptor   *encryption.Encryptor
	validate    *validator.Validate
	devMode     bool
	middlewares []httpapi.Middleware
}

type NewMessagesEndpointParams struct {
	fx.In

	GatewayClient *gateway.Client
	Encryptor     *encryption.Encryptor
	Validate      *validator.Validate
	Config        config.Config
	Authorization httpapi.Middleware
}

type NewMessagesEndpointResult struct {
	fx.Out

	Route httpapi.Route `group:"routes"`
}

func NewMessagesEndpoint(in NewMessagesEndpointParams) NewMessagesEndpointResult {
	return NewMessagesEndpointResult{
		Route: &MessagesEndpoint{
			gateway:     in.GatewayClient,
			encryptor:   in.Encryptor,
			validate:    in.Validate,
			devMode:     in.Config.DevMode,
			middlewares: []httpapi.Middleware{in.Authorization},
		},
	}
}

func (e *MessagesEndpoint) Register(mux *http.ServeMux) {
	httpapi.Handle(mux, "POST /messages", http.HandlerFunc(e.enqueueMessage), e.Middlewares())
}

func (e *MessagesEndpoint) Middlewares() []httpapi.Middleware {
	return e.middlewares
}

func (e *MessagesEndpoint) enqueueMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req httpapi.SendSMSRequest
	if err := httpapi.DecodeJSONStrict(r, &req); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	if err := e.validate.Struct(req); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	payload, err := e.buildGatewayPayload(req)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to prepare request: " + err.Error()})
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to encode request: " + err.Error()})
		return
	}

	gatewayRequest := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Path: messagesPath, RawQuery: r.URL.RawQuery},
		Header: make(http.Header),
		Body:   io.NopCloser(bytes.NewReader(body)),
	}
	gatewayRequest.Header.Set("Content-Type", "application/json")

	response, err := e.gateway.Forward(r.Context(), gatewayRequest)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httpapi.WriteGatewayResponse(w, response)
}

func (e *MessagesEndpoint) buildGatewayPayload(req httpapi.SendSMSRequest) (httpapi.GatewaySendRequest, error) {
	phones := make([]string, 0, len(req.PhoneNumbers))
	for _, number := range req.PhoneNumbers {
		if e.devMode {
			phones = append(phones, number)
			continue
		}

		encrypted, err := e.encryptor.Encrypt(number)
		if err != nil {
			return httpapi.GatewaySendRequest{}, fmt.Errorf("encrypt phone number: %w", err)
		}
		phones = append(phones, encrypted)
	}

	payload := httpapi.GatewaySendRequest{
		DeviceID:           req.DeviceID,
		ID:                 req.ID,
		IsEncrypted:        !e.devMode,
		PhoneNumbers:       phones,
		Priority:           req.Priority,
		ScheduleAt:         req.ScheduleAt,
		SimNumber:          req.SimNumber,
		TTL:                req.TTL,
		ValidUntil:         req.ValidUntil,
		WithDeliveryReport: req.WithDeliveryReport,
	}

	if req.DataMessage != nil {
		if e.devMode {
			payload.DataMessage = &httpapi.DataMessage{Data: req.DataMessage.Data, Port: req.DataMessage.Port}
			return payload, nil
		}

		encryptedData, err := e.encryptor.Encrypt(req.DataMessage.Data)
		if err != nil {
			return httpapi.GatewaySendRequest{}, fmt.Errorf("encrypt dataMessage.data: %w", err)
		}
		payload.DataMessage = &httpapi.DataMessage{Data: encryptedData, Port: req.DataMessage.Port}
		return payload, nil
	}

	text := req.Message
	if req.TextMessage != nil {
		text = req.TextMessage.Text
	}

	if e.devMode {
		payload.TextMessage = &httpapi.TextMessage{Text: text}
		return payload, nil
	}

	encryptedText, err := e.encryptor.Encrypt(text)
	if err != nil {
		return httpapi.GatewaySendRequest{}, fmt.Errorf("encrypt text message: %w", err)
	}

	payload.TextMessage = &httpapi.TextMessage{Text: encryptedText}
	return payload, nil
}
