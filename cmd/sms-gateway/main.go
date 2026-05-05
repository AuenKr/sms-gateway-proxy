package main

import (
	"log"

	"sms-gateway/internal"
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
		fx.Invoke(httpapi.RegisterRoutes, httpapi.RegisterServer),
	)

	app.Run()
	if err := app.Err(); err != nil {
		log.Fatal(err)
	}
}
