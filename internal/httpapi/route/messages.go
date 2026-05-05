package route

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"sms-gateway/internal/encryption"
	"sms-gateway/internal/gateway"
	"sms-gateway/internal/httpapi"
)

const messagesPath = "/3rdparty/v1/messages"

type MessagesEndpoint struct {
	gateway     *gateway.Client
	encryptor   *encryption.Encryptor
	middlewares []httpapi.Middleware
}

func NewMessagesEndpoint(gatewayClient *gateway.Client, encryptor *encryption.Encryptor, authorization httpapi.Middleware) httpapi.Route {
	return &MessagesEndpoint{gateway: gatewayClient, encryptor: encryptor, middlewares: []httpapi.Middleware{authorization}}
}

func (e *MessagesEndpoint) Register(mux *http.ServeMux) {
	httpapi.Handle(mux, "GET /messages", http.HandlerFunc(e.listMessages), e.Middlewares())
	httpapi.Handle(mux, "GET /messages/{id}", http.HandlerFunc(e.getMessageByID), e.Middlewares())
	httpapi.Handle(mux, "POST /messages", http.HandlerFunc(e.enqueueMessage), e.Middlewares())
}

func (e *MessagesEndpoint) Middlewares() []httpapi.Middleware {
	return e.middlewares
}

func (e *MessagesEndpoint) listMessages(w http.ResponseWriter, r *http.Request) {
	response, err := e.gateway.Do(r.Context(), http.MethodGet, messagesPath, r.URL.Query(), nil)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	httpapi.WriteGatewayResponse(w, response)
}

func (e *MessagesEndpoint) getMessageByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		httpapi.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "message id is required"})
		return
	}

	response, err := e.gateway.Do(r.Context(), http.MethodGet, messagesPath+"/"+id, r.URL.Query(), nil)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	httpapi.WriteGatewayResponse(w, response)
}

func (e *MessagesEndpoint) enqueueMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req httpapi.SendSMSRequest
	if err := httpapi.DecodeJSONStrict(r, &req); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	if err := validateSendSMSRequest(req); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	payload, err := e.buildGatewayPayload(req)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to encrypt request: " + err.Error()})
		return
	}

	response, err := e.gateway.DoJSON(r.Context(), http.MethodPost, messagesPath, r.URL.Query(), payload)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httpapi.WriteGatewayResponse(w, response)
}

func validateSendSMSRequest(req httpapi.SendSMSRequest) error {
	if len(req.PhoneNumbers) == 0 {
		return errors.New("phoneNumbers is required")
	}

	for _, number := range req.PhoneNumbers {
		if strings.TrimSpace(number) == "" {
			return errors.New("phoneNumbers cannot contain empty values")
		}
	}

	messageKinds := 0
	if strings.TrimSpace(req.Message) != "" {
		messageKinds++
	}
	if req.TextMessage != nil {
		messageKinds++
		if strings.TrimSpace(req.TextMessage.Text) == "" {
			return errors.New("textMessage.text is required when textMessage is provided")
		}
	}
	if req.DataMessage != nil {
		messageKinds++
		if strings.TrimSpace(req.DataMessage.Data) == "" {
			return errors.New("dataMessage.data is required when dataMessage is provided")
		}
		if req.DataMessage.Port < 0 || req.DataMessage.Port > 65535 {
			return errors.New("dataMessage.port must be between 0 and 65535")
		}
	}

	if messageKinds == 0 {
		return errors.New("exactly one of message, textMessage, or dataMessage is required")
	}
	if messageKinds > 1 {
		return errors.New("only one of message, textMessage, or dataMessage may be provided")
	}

	if req.TTL != nil && req.ValidUntil != "" {
		return errors.New("ttl and validUntil are mutually exclusive")
	}

	if req.ScheduleAt != "" {
		if _, err := time.Parse(time.RFC3339, req.ScheduleAt); err != nil {
			return errors.New("scheduleAt must be RFC3339 format")
		}
	}

	if req.ValidUntil != "" {
		if _, err := time.Parse(time.RFC3339, req.ValidUntil); err != nil {
			return errors.New("validUntil must be RFC3339 format")
		}
	}

	return nil
}

func (e *MessagesEndpoint) buildGatewayPayload(req httpapi.SendSMSRequest) (httpapi.GatewaySendRequest, error) {
	encryptedPhones := make([]string, 0, len(req.PhoneNumbers))
	for _, number := range req.PhoneNumbers {
		encrypted, err := e.encryptor.Encrypt(number)
		if err != nil {
			return httpapi.GatewaySendRequest{}, fmt.Errorf("encrypt phone number: %w", err)
		}
		encryptedPhones = append(encryptedPhones, encrypted)
	}

	payload := httpapi.GatewaySendRequest{
		DeviceID:           req.DeviceID,
		ID:                 req.ID,
		IsEncrypted:        true,
		PhoneNumbers:       encryptedPhones,
		Priority:           req.Priority,
		ScheduleAt:         req.ScheduleAt,
		SimNumber:          req.SimNumber,
		TTL:                req.TTL,
		ValidUntil:         req.ValidUntil,
		WithDeliveryReport: req.WithDeliveryReport,
	}

	if req.DataMessage != nil {
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

	encryptedText, err := e.encryptor.Encrypt(text)
	if err != nil {
		return httpapi.GatewaySendRequest{}, fmt.Errorf("encrypt text message: %w", err)
	}

	payload.TextMessage = &httpapi.TextMessage{Text: encryptedText}
	return payload, nil
}
