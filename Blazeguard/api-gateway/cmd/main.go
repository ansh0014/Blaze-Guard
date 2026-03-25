package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	h := handler.NewGatewayHandler(cfg)
	defer h.Close()
	h.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      middleware.Chain(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("API gateway starting on :%s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}
