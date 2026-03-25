package handler

import (
	"api-gateway/internal/config"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type GatewayHandler struct {
	writers map[string]*kafka.Writer
}

func NewGatewayHandler(cfg config.Config) *GatewayHandler {
	return &GatewayHandler{
		writers: map[string]*kafka.Writer{
			"fire_detected": kafka.NewWriter(kafka.WriterConfig{
				Brokers:      []string{cfg.KafkaBroker},
				Topic:        "fire_detected",
				RequiredAcks: int(kafka.RequireOne),
				Async:        false,
			}),
			"fire_prevention_check": kafka.NewWriter(kafka.WriterConfig{
				Brokers:      []string{cfg.KafkaBroker},
				Topic:        "fire_prevention_check",
				RequiredAcks: int(kafka.RequireOne),
				Async:        false,
			}),
		},
	}
}
func (h *GatewayHandler) Close() {
	for _, w := range h.writers {
		_ = w.Close()
	}
}
func (h *GatewayHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/api/v1/events/fire-detected", h.PublishFireDetected)
	mux.HandleFunc("/api/v1/events/fire-prevention-check", h.PublishFirePreventionCheck)
}
func (h *GatewayHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"service": "api-gateway",
		"status":  "ok",
	})
}
func (h *GatewayHandler) PublishFireDetected(w http.ResponseWriter, r *http.Request) {
	h.publishEvent(w, r, "fire_detected")
}
func (h *GatewayHandler) PublishFirePreventionCheck(w http.ResponseWriter, r *http.Request) {
	h.publishEvent(w, r, "fire_prevention_check")
}
func (h *GatewayHandler) publishEvent(w http.ResponseWriter, r *http.Request, topic string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if _, ok := payload["timestamp"]; !ok {
		payload["timestamp"] = time.Now().Format(time.RFC3339)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "encode failed", http.StatusInternalServerError)
		return
	}

	writer, ok := h.writers[topic]
	if !ok {
		http.Error(w, "topic not configured", http.StatusInternalServerError)
		return
	}

	if err := writer.WriteMessages(context.Background(), kafka.Message{Value: data}); err != nil {
		http.Error(w, "kafka failed to publish event", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"status": "queued",
		"topic":  topic,
	})
}
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
