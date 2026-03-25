package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	grpcserver "orchestrator/grpc"
	"orchestrator/monitor"
	"orchestrator/registry"
	"orchestrator/router"
	"orchestrator/state"
)

func main() {
	reg := registry.NewAgentRegistry()
	reg.BootstrapDefaults()

	store := state.NewInMemoryStore()
	msgRouter := router.NewMessageRouter(reg, store)

	healthMonitor := monitor.NewHealthMonitor(reg, 15*time.Second)
	go healthMonitor.Start()

	grpcPort := os.Getenv("ORCHESTRATOR_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}
	go func() {
		log.Printf("[Orchestrator] gRPC listening on :%s", grpcPort)
		if err := grpcserver.Start(grpcPort, reg, msgRouter); err != nil {
			log.Fatal(err)
		}
	}()

	mux := &http.ServeMux{}
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service": "orchestrator",
			"status":  "ok",
		})
	})
	mux.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(reg.List())
	})
	mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var msg router.A2AMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		if err := msgRouter.Forward(r.Context(), msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "routed",
			"to":     msg.To,
		})
	})

	port := os.Getenv("ORCHESTRATOR_PORT")
	if port == "" {
		port = "8090"
	}

	log.Printf("[Orchestrator] HTTP listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
