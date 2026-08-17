package main

import (
	"blazeguard/shared"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	s2 "blazeguard/agent/prediction/server"

	"github.com/segmentio/kafka-go"
)

type SeverityRequest struct {
	Brightness        float64 `json:"brightness"`
	BrightT31         float64 `json:"bright_t31"`
	Scan              float64 `json:"scan"`
	Track             float64 `json:"track"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Confidence        float64 `json:"confidence"`
	WindSpeed         float64 `json:"wind_speed"`
	WindDeg           float64 `json:"wind_deg"`
	VegetationDensity float64 `json:"vegetation_density"`
}

type SeverityResponse struct {
	Severity     string    `json:"severity"`
	RadiusMeters int       `json:"radius_meters"`
	Corridors    []string  `json:"corridors"`
	Location     []float64 `json:"location"`
}

type RiskRequest struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
}

type RiskResponse struct {
	RiskScore float64 `json:"risk_score"`
	Level     string  `json:"level"`
}

func main() {
	shared.LoadEnv()

	if err := shared.RequireEnv("KAFKA_BROKER", "EVENT_VERSION"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("[Prediction Agent] Starting...")

	s2.SetMessageHandler(handleA2AMessage)
	go consumeTopic("wheather_fire_predictions", handleWeatherPrediction)
	go consumeTopic("fire_detected", handleConfirmedFire)

	go s2.StartHTTPServer()

	select {}
}

func consumeTopic(topic string, handler func([]byte)) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{shared.GetEnv("KAFKA_BROKER", "localhost:9092")},
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
		return
	}
	zoneID, _ := event["zone_id"].(string)
	riskScore := calculateRiskScore(event)
	fmt.Printf("[Prediction] Weather-based fire risk | Zone: %s | Score: %.2f\n",
		zoneID, riskScore)
	if riskScore > 0.5 {
		fmt.Printf("[Prediction] High risk detected in zone %s - triggering prevention\n", zoneID)
		triggerPreventionAction(event, riskScore)
		
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

	lat := toFloat(payload["latitude"])
	lon := toFloat(payload["longitude"])
	conf := toFloat(payload["confidence"])
	if conf == 0 {
		detection, _ := payload["detection"].(map[string]interface{})
		if detection != nil {
			conf = toFloat(detection["confidence"])
		}
	}
	if conf == 0 {
		conf = 0.90
	}

	brightness := toFloat(payload["brightness"])
	if brightness == 0 {
		brightness = 310.0
	}
	brightT31 := toFloat(payload["bright_t31"])
	if brightT31 == 0 {
		brightT31 = 295.0
	}
	scan := toFloat(payload["scan"])
	if scan == 0 {
		scan = 0.5
	}
	track := toFloat(payload["track"])
	if track == 0 {
		track = 0.5
	}

	var windSpeed, windDeg, vegDensity float64
	env, ok := payload["environment"].(map[string]interface{})
	if ok {
		windSpeed = toFloat(env["wind_speed"])
		windDeg = toFloat(env["wind_deg"])
		if windDeg == 0 {
			windDeg = toFloat(env["wind_degree"])
		}
		vegDensity = toFloat(env["vegetation_density"])
	}
	if windSpeed == 0 {
		windSpeed = 18.0
	}
	if windDeg == 0 {
		windDeg = 180.0
	}
	if vegDensity == 0 {
		vegDensity = 0.7
	}

	mlResult, err := callMLModelPredictSeverity(SeverityRequest{
		Brightness:        brightness,
		BrightT31:         brightT31,
		Scan:              scan,
		Track:             track,
		Latitude:          lat,
		Longitude:         lon,
		Confidence:        conf,
		WindSpeed:         windSpeed,
		WindDeg:           windDeg,
		VegetationDensity: vegDensity,
	})

	var corridors []string
	if err != nil {
		log.Printf("[Prediction] ML prediction call failed, using fallback: %v", err)
		corridors = []string{"corridor_north", "corridor_east"}
	} else {
		fmt.Printf("[Prediction] ML Severity: %s | Spread Radius: %dm | Predicted Corridors: %v\n", 
			mlResult.Severity, mlResult.RadiusMeters, mlResult.Corridors)
		corridors = mlResult.Corridors
	}

	go shared.SendToAgent("SELF_EVOLVING_AGENT", shared.A2AMessage{
		From:      "prediction_agent",
		To:        "self_evolving_agent",
		EventType: "FIRE_SPREAD_PREDICTION",
		Payload: map[string]interface{}{
			"zone_id":             zoneID,
			"predicted_corridors": corridors,
			"environment":         payload["environment"],
			"timestamp":           payload["timestamp"],
			"model_version":       "v1.1",
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

	lat := toFloat(event["latitude"])
	lon := toFloat(event["longitude"])
	conf := toFloat(event["confidence"])
	if conf == 0 {
		detection, _ := event["detection"].(map[string]interface{})
		if detection != nil {
			conf = toFloat(detection["confidence"])
		}
	}
	if conf == 0 {
		conf = 0.90
	}

	brightness := toFloat(event["brightness"])
	if brightness == 0 {
		brightness = 310.0
	}
	brightT31 := toFloat(event["bright_t31"])
	if brightT31 == 0 {
		brightT31 = 295.0
	}
	scan := toFloat(event["scan"])
	if scan == 0 {
		scan = 0.5
	}
	track := toFloat(event["track"])
	if track == 0 {
		track = 0.5
	}

	var windSpeed, windDeg, vegDensity float64
	env, ok := event["environment"].(map[string]interface{})
	if ok {
		windSpeed = toFloat(env["wind_speed"])
		windDeg = toFloat(env["wind_deg"])
		if windDeg == 0 {
			windDeg = toFloat(env["wind_degree"])
		}
		vegDensity = toFloat(env["vegetation_density"])
	}
	if windSpeed == 0 {
		windSpeed = 18.0
	}
	if windDeg == 0 {
		windDeg = 180.0
	}
	if vegDensity == 0 {
		vegDensity = 0.7
	}

	mlResult, err := callMLModelPredictSeverity(SeverityRequest{
		Brightness:        brightness,
		BrightT31:         brightT31,
		Scan:              scan,
		Track:             track,
		Latitude:          lat,
		Longitude:         lon,
		Confidence:        conf,
		WindSpeed:         windSpeed,
		WindDeg:           windDeg,
		VegetationDensity: vegDensity,
	})

	var corridors []string
	if err != nil {
		log.Printf("[Prediction] ML prediction call failed, using fallback: %v", err)
		corridors = []string{"corridor_north", "corridor_east"}
	} else {
		fmt.Printf("[Prediction] ML Severity: %s | Spread Radius: %dm | Predicted Corridors: %v\n", 
			mlResult.Severity, mlResult.RadiusMeters, mlResult.Corridors)
		corridors = mlResult.Corridors
	}

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
	var temp, humidity, windSpeed float64
	env, ok := event["environment"].(map[string]interface{})
	if !ok {
		env, ok = event["weather"].(map[string]interface{})
	}
	
	if ok {
		temp = toFloat(env["temperature"])
		humidity = toFloat(env["humidity"])
		windSpeed = toFloat(env["wind_speed"])
	} else {
		temp = toFloat(event["temperature"])
		humidity = toFloat(event["humidity"])
		windSpeed = toFloat(event["wind_speed"])
	}

	if temp == 0 {
		temp = 25.0
	}
	if humidity == 0 {
		humidity = 50.0
	}
	if windSpeed == 0 {
		windSpeed = 15.0
	}

	mlResult, err := callMLModelPredictRisk(RiskRequest{
		Temperature: temp,
		Humidity:    humidity,
		WindSpeed:   windSpeed,
	})

	if err != nil {
		log.Printf("[Prediction] ML predict failed for risk score, using fallback formula: %v", err)
		return (1.0 - humidity/100.0)*0.4 + (windSpeed/50.0)*0.3 + (temp/50.0)*0.3
	}

	fmt.Printf("[Prediction] ML Risk Score: %.2f | Level: %s\n", mlResult.RiskScore, mlResult.Level)
	return mlResult.RiskScore
}

func callMLModelPredictSeverity(payload SeverityRequest) (*SeverityResponse, error) {
	mlURL := shared.GetEnv("ML_MODEL_URL", "http://localhost:9000")
	url := mlURL + "/predict-severity"

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ML model returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result SeverityResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func callMLModelPredictRisk(payload RiskRequest) (*RiskResponse, error) {
	mlURL := shared.GetEnv("ML_MODEL_URL", "http://localhost:9000")
	url := mlURL + "/predict-risk"

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ML model returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result RiskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func triggerPreventionAction(event map[string]interface{}, score float64) {
	zoneID, _ := event["zone_id"].(string)
	fmt.Printf("[Prediction] Triggering prevention actions for zone %s with risk score %.2f\n", zoneID, score)

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
