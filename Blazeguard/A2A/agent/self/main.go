package main

import (
	s5 "blazeguard/agent/self/server"
	"blazeguard/shared"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type SelfState struct {
	mu sync.RWMutex

	DetectionsTotal          int64
	DetectionsLowConfidence  int64
	PredictionsTotal         int64
	PredictionsLowConfidence int64
	LogisticsRoutesTotal     int64
	LogisticsRoutesFailed    int64
	CitizenAlertsTotal       int64
	CitizenAlertsCritical    int64
	LastUpdatedUnix          int64
}

var state = &SelfState{}

func main() {
	shared.LoadEnv()

	if err := shared.RequireEnv("KAFKA_BROKER", "EVENT_VERSION"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("[Self Agent] Starting...")
	s5.SetMessageHandler(handleA2AMessage)
	go consumeTopic("fire_detected", handleDetectionEvent)
	go consumeTopic("fire_prevention_check", handlePredictionEvent)
	go consumeTopic("logistics_routes", handleLogisticsEvent)
	go consumeTopic("citizen_alerts", handleCitizenAlertEvent)
	go reportLoop()
	go s5.StartHTTPServer()
	select {}
}
func consumeTopic(topic string, handler func([]byte)) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{shared.GetEnv("KAFKA_BROKER", "localhost:9092")},
		Topic:    topic,
		GroupID:  "self_agent_group_" + topic,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()
	fmt.Printf("[Self Agent] Listening to %s\n", topic)
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[Self Agent] Kafka read error (%s): %v", topic, err)
			continue
		}
		handler(msg.Value)
	}
}
func handleA2AMessage(eventType string, payload map[string]interface{}) {
	switch eventType {
	case "MODEL_FEEDBACK":
		handleModelFeedback(payload)
	case "MODEL_EVAL_RESULT":
		handleModelEvalResult(payload)
	default:
		fmt.Printf("[Self Agent] Unknown event type: %s\n", eventType)
	}
}
func handleDetectionEvent(data []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		return
	}
	confidence := toFloat(event["confidence"])
	state.mu.Lock()
	state.DetectionsTotal++
	if confidence > 0 && confidence < 0.55 {
		state.DetectionsLowConfidence++
	}
	state.LastUpdatedUnix = time.Now().Unix()
	state.mu.Unlock()
}
func handlePredictionEvent(data []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		return
	}
	confidence := toFloat(event["confidence"])
	state.mu.Lock()
	state.PredictionsTotal++
	if confidence > 0 && confidence < 0.60 {
		state.PredictionsLowConfidence++
	}
	state.LastUpdatedUnix = time.Now().Unix()
	state.mu.Unlock()
}
func handleLogisticsEvent(data []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		return
	}
	route, _ := event["route"].(map[string]interface{})
	stationID, _ := route["station_id"].(string)
	state.mu.Lock()
	state.LogisticsRoutesTotal++
	if stationID == "" {
		state.LogisticsRoutesFailed++
	}
	state.LastUpdatedUnix = time.Now().Unix()
	state.mu.Unlock()
}
func handleCitizenAlertEvent(data []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		return
	}
	severity, _ := event["severity"].(string)
	state.mu.Lock()
	state.CitizenAlertsTotal++
	if severity == "CRITICAL" {
		state.CitizenAlertsCritical++
	}
	state.LastUpdatedUnix = time.Now().Unix()
	state.mu.Unlock()
}
func handleModelFeedback(payload map[string]interface{}) {
	agent, _ := payload["agent"].(string)
	score := toFloat(payload["score"])
	fmt.Printf("[Self Agent] Model feedback | agent=%s score=%.4f\n", agent, score)
}
func handleModelEvalResult(payload map[string]interface{}) {
	agent, _ := payload["agent"].(string)
	version, _ := payload["version"].(string)
	pass, _ := payload["pass"].(bool)
	fmt.Printf("[Self Agent] Eval result | agent=%s version=%s pass=%v\n", agent, version, pass)
}
func reportLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		report := buildReport()
		publishToKafka("self_agent_reports", report)

		go shared.SendToAgent("PREDICTION_AGENT", shared.A2AMessage{
			From:      "self_agent",
			To:        "prediction_agent",
			EventType: "SELF_AGENT_REPORT",
			Payload:   report,
		})
	}
}
func buildReport() map[string]interface{} {
	state.mu.RLock()
	defer state.mu.RUnlock()

	return map[string]interface{}{
		"agent": "self_agent",
		"metrics": map[string]interface{}{
			"detections_total":           state.DetectionsTotal,
			"detections_low_confidence":  state.DetectionsLowConfidence,
			"predictions_total":          state.PredictionsTotal,
			"predictions_low_confidence": state.PredictionsLowConfidence,
			"logistics_routes_total":     state.LogisticsRoutesTotal,
			"logistics_routes_failed":    state.LogisticsRoutesFailed,
			"citizen_alerts_total":       state.CitizenAlertsTotal,
			"citizen_alerts_critical":    state.CitizenAlertsCritical,
		},
		"last_updated_unix": state.LastUpdatedUnix,
		"generated_at_unix": time.Now().Unix(),
	}
}
func publishToKafka(topic string, payload map[string]interface{}) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   topic,
	})
	defer writer.Close()
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	err = writer.WriteMessages(context.Background(), kafka.Message{Value: data})
	if err != nil {
		log.Printf("[Self Agent] Kafka publish error (%s): %v", topic, err)
	}
}
func toFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int64:
		return float64(t)
	default:
		return 0
	}
}
