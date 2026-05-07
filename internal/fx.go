package internal

import (
	"sms-gateway/internal/config"
	"sms-gateway/internal/encryption"
	"sms-gateway/internal/gateway"
	"sms-gateway/internal/httpapi"
	"sms-gateway/internal/httpapi/middleware"
	"sms-gateway/internal/httpapi/route"

	"go.uber.org/fx"
)

var Module = fx.Provide(
	config.Load,
	encryption.New,
	gateway.NewHTTPClient,
	gateway.NewClient,
	httpapi.NewServeMux,
	httpapi.NewHandler,
	httpapi.NewValidator,
	middleware.NewLogger,
	middleware.NewAuthorization,
	middleware.NewContentLengthValidation,
	middleware.NewRequestLogging,
	route.NewHealthEndpoint,
	route.NewMessagesEndpoint,
	route.NewProxyEndpoint,
)
