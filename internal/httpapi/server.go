package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"sort"

	"go.uber.org/fx"

	"sms-gateway/internal/config"
	"sms-gateway/internal/gateway"
)

type Route interface {
	Register(*http.ServeMux)
	Middlewares() []Middleware
}

type Middleware func(http.Handler) http.Handler

type OrderedMiddleware struct {
	Order      int
	Middleware Middleware
}

func NewServeMux() *http.ServeMux {
	return http.NewServeMux()
}

type middlewareParams struct {
	fx.In

	Mux         *http.ServeMux
	Middlewares []OrderedMiddleware `group:"middlewares"`
}

func NewHandler(params middlewareParams) http.Handler {
	sort.SliceStable(params.Middlewares, func(i, j int) bool {
		return params.Middlewares[i].Order < params.Middlewares[j].Order
	})

	middlewares := make([]Middleware, 0, len(params.Middlewares))
	for _, entry := range params.Middlewares {
		middlewares = append(middlewares, entry.Middleware)
	}

	return ApplyMiddlewares(params.Mux, middlewares)
}

func ApplyMiddlewares(handler http.Handler, middlewares []Middleware) http.Handler {
	wrapped := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}

func Handle(mux *http.ServeMux, pattern string, handler http.Handler, middlewares []Middleware) {
	mux.Handle(pattern, ApplyMiddlewares(handler, middlewares))
}

type routeParams struct {
	fx.In

	Mux    *http.ServeMux
	Routes []Route `group:"routes"`
}

func RegisterRoutes(params routeParams) {
	for _, route := range params.Routes {
		route.Register(params.Mux)
	}
}

func RegisterServer(lifecycle fx.Lifecycle, cfg config.Config, handler http.Handler) {
	var listener net.Listener

	server := &http.Server{Addr: ":" + cfg.Port, Handler: handler}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			listener, err = net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}

			go func() {
				if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
					_ = listener.Close()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

func decodeJSONStrict(r *http.Request, dst any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("request body must contain a single JSON object")
	}
	return nil
}

func DecodeJSONStrict(r *http.Request, dst any) error {
	return decodeJSONStrict(r, dst)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	writeJSON(w, status, payload)
}

func writeGatewayResponse(w http.ResponseWriter, response gateway.Response) {
	contentType := response.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(response.StatusCode)
	_, _ = w.Write(response.Body)
}

func WriteGatewayResponse(w http.ResponseWriter, response gateway.Response) {
	writeGatewayResponse(w, response)
}
