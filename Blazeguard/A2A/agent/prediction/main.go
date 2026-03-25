package main

import (
	"blazeguard/shared"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"os"

	s2 "blazeguard/agent/prediction/server"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func main() {
	godotenv.Load("../../.env")

	fmt.Println("[Prediction Agent] Starting...")

	s2.SetMessageHandler(handleA2AMessage)
	go consumeTopic("wheather_fire_predictions", handleWeatherPrediction)
	go consumeTopic("fire_detected", handleConfirmedFire)

	go s2.StartHTTPServer()

	select {}
}
func consumeTopic(topic string, handler func([]byte)) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   topic,
		GroupID: "prediction_service_group" + topic,
	})
	defer reader.Close()
	fmt.Printf("[Prediction Agent] Listening to %s\n", topic)
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[Prediction] Kafka error: %v", err)
			continue
		}
		handler(msg.Value)
	}
}

func handleA2AMessage(eventType string, payload map[string]interface{}) {
	switch eventType {
	case "FIRE_DETECTED":
		handleA2AFireDetected(payload)
	case "NO_FIRE_PREVENTION_CHECK":
		handleA2APreventionCheck(payload)
	default:
		fmt.Printf("[Prediction] Unknown event type: %s\n", eventType)
	}
}
func handleWeatherPrediction(data []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[Prediction] Json parse error : %v", err)
	}
	zoneID, _ := event["zone_id"].(string)
	riskScore := calculateRiskScore(event)
	fmt.Printf("[Prediction] Weather-based fire risk | Zone: %s | Score: %.2f\n",
		zoneID, riskScore)
	if riskScore > 0.5 {
		// i have to change the value when the ml model complete
		fmt.Printf("[Prediction] High risk detected in zone %s -triggering prevention\n", zoneID)
		triggerPreventionAction(event, riskScore)
		// Send to self-evolving for learning
		go shared.SendToAgent("SELF_EVOLVING_AGENT", shared.A2AMessage{
			From:      "prediction_agent",
			To:        "self_evolving_agent",
			EventType: "FIRE_RISK_PREDICTION",
			Payload: map[string]interface{}{
				"zone_id":      zoneID,
				"risk_score":   riskScore,
				"weather_data": event["weather"],
				"timestamp":    event["timestamp"],
			},
		})

	}

}
func handleA2AFireDetected(payload map[string]interface{}) {
	zoneID, _ := payload["zone_id"].(string)
	fmt.Printf("[Prediction] A2A Fire detected in zone %s - Running spread analysis\n", zoneID)

	corridors := runSpreadModel(payload)
	go shared.SendToAgent("SELF_EVOLVING_AGENT", shared.A2AMessage{
		From:      "prediction_agent",
		To:        "self_evolving_agent",
		EventType: "FIRE_SPREAD_PREDICTION",
		Payload: map[string]interface{}{
			"zone_id":             zoneID,
			"predicted_corridors": corridors,
			"environment":         payload["environment"],
			"timestamp":           payload["timestamp"],
			"model_version":       "v1.0",
		},
	})

	go shared.SendToAgent("LOGISTICS_AGENT", shared.A2AMessage{
		From:      "prediction_agent",
		To:        "logistics_agent",
		EventType: "FIRE_SPREAD_PREDICTION",
		Payload: map[string]interface{}{
			"zone_id":   zoneID,
			"corridors": corridors,
			"event":     payload,
		},
	})

	go shared.SendToAgent("CITIZEN_ALERT_AGENT", shared.A2AMessage{
		From:      "prediction_agent",
		To:        "citizen_alert_agent",
		EventType: "EVACUATION_CORRIDORS",
		Payload: map[string]interface{}{
			"zone_id":   zoneID,
			"corridors": corridors,
		},
	})

	fmt.Printf("[Prediction] Sent spread prediction to Logistics, Citizen Alert, and Self-Evolving\n")
}

func handleA2APreventionCheck(payload map[string]interface{}) {
	zoneID, _ := payload["zone_id"].(string)
	fmt.Printf("[Prediction] A2A Prevention check for zone %s\n", zoneID)

	riskScore := calculateRiskScore(payload)
	fmt.Printf("[Prediction] Prevention risk score: %.2f | zone: %s\n", riskScore, zoneID)

	if riskScore > 0.5 {
		triggerPreventionAction(payload, riskScore)
	}
}

func handleConfirmedFire(data []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[Prediction] JSON parse error: %v", err)
		return
	}

	zoneID, _ := event["zone_id"].(string)
	fmt.Printf("[Prediction] Kafka confirmed fire in zone %s - Running spread model\n", zoneID)

	corridors := runSpreadModel(event)

	go shared.SendToAgent("SELF_EVOLVING_AGENT", shared.A2AMessage{
		From:      "prediction_agent",
		To:        "self_evolving_agent",
		EventType: "FIRE_SPREAD_PREDICTION",
		Payload: map[string]interface{}{
			"zone_id":             zoneID,
			"predicted_corridors": corridors,
			"event":               event,
		},
	})

	go shared.SendToAgent("CITIZEN_ALERT_AGENT", shared.A2AMessage{
		From:      "prediction_agent",
		To:        "citizen_alert_agent",
		EventType: "EVACUATION_CORRIDORS",
		Payload: map[string]interface{}{
			"zone_id":   zoneID,
			"corridors": corridors,
		},
	})

	fmt.Printf("[Prediction] Spread predicted - %d corridors identified\n", len(corridors))
}
func calculateRiskScore(event map[string]interface{}) float64 {
	env, ok := event["environment"].(map[string]interface{})
	if !ok {
		log.Println("[Prediction] No environment data available")
		return 0.0
	}

	humidity, _ := env["humidity"].(float64)
	windSpeed, _ := env["wind_speed"].(float64)
	temperature, _ := env["temperature"].(float64)

	// Basic risk formula until ML model is ready
	// Low humidity + high wind + high temp = higher risk
	//   check it i have to change the value when the ml model complete
	riskScore := (1.0-humidity/100.0)*0.4 + (windSpeed/50.0)*0.3 + (temperature/50.0)*0.3

	return riskScore
}

func runSpreadModel(event map[string]interface{}) []string {
	// TODO: Integrate ML model for fire spread prediction
	// This will use environmental data, wind direction, terrain, etc.
	// to predict fire spread corridors

	// Placeholder corridors for now
	return []string{"corridor_north", "corridor_east"}
}

func triggerPreventionAction(event map[string]interface{}, score float64) {
	zoneID, _ := event["zone_id"].(string)
	fmt.Printf("[Prediction] Triggering prevention actions for zone %s with risk score %.2f\n", zoneID, score)

	// Send alert to Citizen Alert Agent
	go shared.SendToAgent("CITIZEN_ALERT_AGENT", shared.A2AMessage{
		From:      "prediction_agent",
		To:        "citizen_alert_agent",
		EventType: "PREVENTION_ALERT",
		Payload: map[string]interface{}{
			"zone_id":    zoneID,
			"risk_score": score,
			"message":    "High fire risk detected. Please follow prevention guidelines.",
		},
	})

	go shared.SendToAgent("LOGISTICS_AGENT", shared.A2AMessage{
		From:      "prediction_agent",
		To:        "logistics_agent",
		EventType: "PREPOSITION_RESOURCES",
		Payload: map[string]interface{}{
			"zone_id":    zoneID,
			"risk_score": score,
			"priority":   "medium",
		},
	})

	fmt.Printf("[Prediction] Prevention actions triggered for zone %s\n", zoneID)
}
