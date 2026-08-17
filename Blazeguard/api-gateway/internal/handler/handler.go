package handler

import (
	"api-gateway/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type GatewayHandler struct {
	writers    map[string]*kafka.Writer
	mlModelURL string
}

func NewGatewayHandler(cfg config.Config) *GatewayHandler {
	return &GatewayHandler{
		mlModelURL: cfg.MLModelURL,
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
			"yolo_fire_events": kafka.NewWriter(kafka.WriterConfig{
				Brokers:      []string{cfg.KafkaBroker},
				Topic:        "yolo_fire_events",
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
	mux.HandleFunc("/api/v1/events/detect-image", h.PublishDetectImage)
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

func (h *GatewayHandler) PublishDetectImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "missing 'image' file parameter", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read optional coordinates
	latVal := r.FormValue("latitude")
	lonVal := r.FormValue("longitude")
	zoneID := r.FormValue("zone_id")
	if zoneID == "" {
		zoneID = "Z-UNKNOWN"
	}

	var lat, lon float64
	if latVal != "" {
		fmt.Sscanf(latVal, "%f", &lat)
	} else {
		lat = 28.6139
	}
	if lonVal != "" {
		fmt.Sscanf(lonVal, "%f", &lon)
	} else {
		lon = 77.2090
	}

	// Read file bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read image file", http.StatusInternalServerError)
		return
	}

	// Call ML model FastAPI
	result, err := h.callMLModelDetect(fileBytes, header.Filename, lat, lon)
	if err != nil {
		http.Error(w, fmt.Sprintf("ML service call failed: %v", err), http.StatusInternalServerError)
		return
	}

	fireDetected, _ := result["fire_detected"].(bool)
	metrics, _ := result["metrics"].(map[string]any)

	// Wrap in A2A-compliant JSON event structure
	payload := map[string]any{
		"event_version": "v1",
		"zone_id":       zoneID,
		"latitude":      lat,
		"longitude":     lon,
		"timestamp":     time.Now().Format(time.RFC3339),
		"detection": map[string]any{
			"fire_detected": fireDetected,
			"confidence":    0.0,
		},
	}

	if fireDetected && metrics != nil {
		if conf, ok := metrics["confidence"].(float64); ok {
			payload["detection"].(map[string]any)["confidence"] = conf
		}
		if brightness, ok := metrics["brightness"].(float64); ok {
			payload["brightness"] = brightness
		}
		if brightT31, ok := metrics["bright_t31"].(float64); ok {
			payload["bright_t31"] = brightT31
		}
		if scan, ok := metrics["scan"].(float64); ok {
			payload["scan"] = scan
		}
		if track, ok := metrics["track"].(float64); ok {
			payload["track"] = track
		}
	}

	// Marshal and publish to Kafka
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "encode failed", http.StatusInternalServerError)
		return
	}

	writer, ok := h.writers["yolo_fire_events"]
	if !ok {
		http.Error(w, "topic not configured", http.StatusInternalServerError)
		return
	}

	if err := writer.WriteMessages(context.Background(), kafka.Message{Value: data}); err != nil {
		http.Error(w, "kafka failed to publish event", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":        "queued",
		"topic":         "yolo_fire_events",
		"fire_detected": fireDetected,
		"payload":       payload,
	})
}

func (h *GatewayHandler) callMLModelDetect(imageBytes []byte, filename string, lat, lon float64) (map[string]any, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, bytes.NewReader(imageBytes)); err != nil {
		return nil, err
	}

	// Add fields
	_ = writer.WriteField("latitude", fmt.Sprintf("%f", lat))
	_ = writer.WriteField("longitude", fmt.Sprintf("%f", lon))

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	url := h.mlModelURL + "/detect"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML model returned status code %d", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
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
