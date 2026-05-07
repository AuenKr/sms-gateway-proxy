package main

import (
	"context"
	"fmt"
	"log"

	"sms-gateway/internal"
	"sms-gateway/internal/config"
	"sms-gateway/internal/httpapi"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}

	app := fx.New(
		internal.Module,
		fx.Invoke(
			httpapi.RegisterRoutes,
			httpapi.RegisterServer,
			func(conf config.Config, lc fx.Lifecycle) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						fmt.Printf("Starting server on port %s\n", conf.Port)
						return nil
					},
					OnStop: func(ctx context.Context) error {
						fmt.Println("Stopping server")
						return nil
					},
				})
			},
		),
	)

	app.Run()
	if err := app.Err(); err != nil {
		log.Fatal(err)
	}
}
