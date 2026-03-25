package main

import (
	"blazeguard/shared"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"os"

	s1 "blazeguard/agent/detection/Server"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func main() {
	godotenv.Load("../../.env")

	fmt.Println("[Detection Agent] Starting Kafka consumer...")
	go consumeKafka()
	go s1.StartHTTPServer()

	select {}
}

func consumeKafka() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		Topic:    "yolo_fire_events",
		GroupID:  "detection-agent-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	fmt.Println("[Detection Agent] Listening to yolo_fire_events...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[Detection] Kafka read error: %v", err)
			continue
		}
		fmt.Println("[Detection] Message received from Kafka")
		processYoloMessage(msg.Value)
	}
}

func processYoloMessage(data []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[Detection] JSON parse error: %v", err)
		return
	}

	detection, ok := event["detection"].(map[string]interface{})
	if !ok {
		log.Println("[Detection] Invalid detection field in event")
		return
	}

	fireDetected, _ := detection["fire_detected"].(bool)
	confidence, _ := detection["confidence"].(float64)
	zoneID, _ := event["zone_id"].(string)

	if fireDetected && confidence > 0.75 { // i have to change the value when the ml model complete
		fmt.Printf("[Detection] FIRE CONFIRMED | Zone: %s | Confidence: %.2f\n", zoneID, confidence)
		go shared.SendToAgent("CITIZEN_ALERT_AGENT", shared.A2AMessage{
			From:      "detection_agent",
			To:        "citizen_alert_agent",
			EventType: "FIRE_DETECTED",
			Payload: map[string]interface{}{
				"zone_id":    zoneID,
				"confidence": confidence,
				"timestamp":  event["timestamp"],
			},
		})
		go shared.SendToAgent("PREDICTION_AGENT", shared.A2AMessage{
			From:      "detection_agent",
			To:        "prediction_agent",
			EventType: "FIRE_DETECTED",
			Payload:   event,
		})

		go shared.SendToAgent("SELF_EVOLVING_AGENT", shared.A2AMessage{
			From:      "detection_agent",
			To:        "self_evolving_agent",
			EventType: "FIRE_DETECTED",
			Payload:   event,
		})
		go shared.SendToAgent("LOGISTICS_AGENT", shared.A2AMessage{
			From:      "detection_agent",
			To:        "logistics_agent",
			EventType: "FIRE_DETECTED",
			Payload:   event,
		})

		go generateAndForwardReport("fire_detected", event)

		publishToKafka("fire_detected", event)

	} else if !fireDetected {
		fmt.Printf("[Detection] No fire|Zone: %s | Running prevention check\n", zoneID)

		go shared.SendToAgent("PREDICTION_AGENT", shared.A2AMessage{
			From:      "detection_agent",
			To:        "prediction_agent",
			EventType: "NO_FIRE_PREVENTION_CHECK",
			Payload:   event,
		})

		go generateAndForwardReport("prevention_alert", event)

	} else {

		fmt.Printf("[Detection] Low confidence %.2f | Zone: %s | Monitoring\n", confidence, zoneID)
	}
}
func generateAndForwardReport(reportType string, event map[string]interface{}) {
	// Generate report text based on event data
	zoneID, _ := event["zone_id"].(string)
	detection, _ := event["detection"].(map[string]interface{})
	confidence, _ := detection["confidence"].(float64)

	var reportText, severity string
	if reportType == "fire_detected" {
		reportText = fmt.Sprintf("Fire detected in zone %s with %.2f%% confidence", zoneID, confidence*100)
		severity = "high"
	} else {
		reportText = fmt.Sprintf("Prevention check recommended for zone %s", zoneID)
		severity = "medium"
	}

	report := shared.Report{
		Report:      reportText,
		Severity:    severity,
		GeneratedAt: time.Now().Format(time.RFC3339),
	}

	fmt.Printf("[Detection] Report generated | type: %s | severity: %s\n",
		reportType, report.Severity)

	publishToKafka("reports", map[string]interface{}{
		"report_type":  reportType,
		"report_text":  report.Report,
		"severity":     report.Severity,
		"generated_at": report.GeneratedAt,
		"zone_id":      event["zone_id"],
		"source":       "detection_agent",
	})
}

func publishToKafka(topic string, payload map[string]interface{}) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   topic,
	})
	defer writer.Close()

	data, _ := json.Marshal(payload)
	err := writer.WriteMessages(context.Background(),
		kafka.Message{Value: data},
	)
	if err != nil {
		log.Printf("[Detection] Kafka publish error on topic %s: %v", topic, err)
		return
	}
	fmt.Printf("[Detection] Published to Kafka topic: %s\n", topic)
}
