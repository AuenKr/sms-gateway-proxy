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
	middleware.NewLogger,
	middleware.NewAuthorization,
	fx.Annotate(middleware.NewContentLengthValidation, fx.ResultTags(`group:"middlewares"`)),
	fx.Annotate(middleware.NewRequestLogging, fx.ResultTags(`group:"middlewares"`)),
	fx.Annotate(route.NewHealthEndpoint, fx.ResultTags(`group:"routes"`)),
	fx.Annotate(route.NewMessagesEndpoint, fx.ResultTags(`group:"routes"`)),
	fx.Annotate(route.NewProxyEndpoint, fx.ResultTags(`group:"routes"`)),
)
