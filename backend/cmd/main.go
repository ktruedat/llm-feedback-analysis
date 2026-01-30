package main

import (
	"log"

	"github.com/ktruedat/llm-feedback-analysis/internal/app"
)

// @title           LLM Feedback Analysis API
// @version         1.0
// @description     This is a LLM Feedback Analysis API server.

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token for authentication. Use the format: "Bearer <token>". This follows RFC 6750 (OAuth 2.0 Bearer Token Usage).

func main() {
	serverApp, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := serverApp.Start(); err != nil {
		log.Fatal(err)
	}
}
