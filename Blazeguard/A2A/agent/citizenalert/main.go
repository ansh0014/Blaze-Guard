package main

// in this code i have to check the spread prediction is implemented in my ai modle
import (
	s4 "blazeguard/agent/citizenalert/server"
	"blazeguard/shared"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	shared.LoadEnv()

	if err := shared.RequireEnv("KAFKA_BROKER", "EVENT_VERSION"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("[Citizen Alert Agent]")
	s4.SetMessageHandler(handleA2AMessage)
	go consumeTopic("fire_detected", handleFireDetected)
	go consumeTopic("fire_prevention_check", handlePreventionCheck)
	go s4.StartHTTPServer()
	select {}
}

func consumeTopic(topic string, handler func([]byte)) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{shared.GetEnv("KAFKA_BROKER", "localhost:9092")},
		Topic:   topic,
		GroupID: "citizen_alert_service_group_" + topic,
	})
	defer reader.Close()
	fmt.Printf("[Citizen Alert Agent] listening to %s\n", topic)
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[Citizen Alert Agent] error reading message from topic %s: %v\n", topic, err)
			continue
		}
		handler(msg.Value)
	}
}

func handleA2AMessage(eventType string, payload map[string]interface{}) {
	switch eventType {
	case "DEPLOYMENT_ROUTE":
		handleDeploymentRoute(payload)
	case "SAFE_ZONES_UPDATE":
		handleSafeZonesUpdate(payload)
	case "EVACUATION_CORRIDORS":
		handleEvacuationCorridors(payload)
	case "PREVENTION_ALERT":
		handlePreventionAlert(payload)
	default:
		fmt.Printf("[Citizen Alert] Unknown event type: %s\n", eventType)
	}
}

func handleFireDetected(data []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[Citizen Alert] JSON Parse error:%v\n", err)
		return
	}
	zoneID, _ := event["zone_id"].(string)
	fmt.Printf("[Citizen Alert] Fire detected in zone %s-sending emergency alert\n", zoneID)
	alert := map[string]interface{}{
		"alert_type": "FIRE_EMERGENCY",
		"severity":   "CRITICAL",
		"zone_id":    zoneID,
		"latitude":   event["latitude"],
		"longitude":  event["longitude"],
		"message":    "EMERGENCY: Wildfire detected in your area. Evacuate immediately!",
		"timestamp":  event["timestamp"],
	}
	sendAlert(alert)
	publishToKafka("citizen_alerts", alert)
	fmt.Printf("[Citizen Alert] Emergency alert sent for zone %s\n", zoneID)
}

func handlePreventionCheck(data []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[Citizen Alert] JSON parse error: %v", err)
		return
	}

	zoneID, _ := event["zone_id"].(string)
	riskScore, _ := event["risk_score"].(float64)

	fmt.Printf("[Citizen Alert] Prevention check for zone %s with risk %.2f\n", zoneID, riskScore)

	if riskScore > 0.7 {
		alert := map[string]interface{}{
			"alert_type": "FIRE_RISK_WARNING",
			"severity":   "HIGH",
			"zone_id":    zoneID,
			"risk_score": riskScore,
			"message":    "High fire risk detected. Stay alert and follow prevention guidelines.",
			"timestamp":  event["timestamp"],
		}

		sendAlert(alert)
		publishToKafka("citizen_alerts", alert)

		fmt.Printf("[Citizen Alert] High risk warning sent for zone %s\n", zoneID)
	}
}

func handleDeploymentRoute(payload map[string]interface{}) {
	zoneID, _ := payload["zone_id"].(string)
	route, _ := payload["route"].(map[string]interface{})
	deployment, _ := payload["deployment"].(map[string]interface{})
	etaMinutes, _ := payload["eta_minutes"].(int64)

	fmt.Printf("[Citizen Alert] Deployment alert for zone %s - ETA: %d min\n", zoneID, etaMinutes)

	alert := map[string]interface{}{
		"alert_type": "HELP_INCOMING",
		"severity":   "INFO",
		"zone_id":    zoneID,
		"message":    fmt.Sprintf("Emergency services en route. ETA: %d minutes", etaMinutes),
		"route":      route,
		"deployment": deployment,
		"timestamp":  getCurrentTimestamp(),
	}

	sendAlert(alert)
	publishToKafka("citizen_alerts", alert)

	fmt.Printf("[Citizen Alert] Deployment notification sent for zone %s\n", zoneID)
}

func handleSafeZonesUpdate(payload map[string]interface{}) {
	zoneID, _ := payload["zone_id"].(string)
	safeZones, _ := payload["safe_zones"].([]interface{})
	evacuationRoutes, _ := payload["evacuation_routes"].([]interface{})

	fmt.Printf("[Citizen Alert] Safe zones update for zone %s - %d zones available\n", zoneID, len(safeZones))

	alert := map[string]interface{}{
		"alert_type":        "SAFE_ZONES",
		"severity":          "WARNING",
		"zone_id":           zoneID,
		"message":           "Evacuate to nearest safe zone immediately!",
		"safe_zones":        safeZones,
		"evacuation_routes": evacuationRoutes,
		"timestamp":         getCurrentTimestamp(),
	}

	sendAlert(alert)
	publishToKafka("citizen_alerts", alert)

	fmt.Printf("[Citizen Alert] Safe zone information sent for zone %s\n", zoneID)
}

func handleEvacuationCorridors(payload map[string]interface{}) {
	zoneID, _ := payload["zone_id"].(string)
	corridors, _ := payload["corridors"].([]interface{})

	fmt.Printf("[Citizen Alert] Evacuation corridors for zone %s - %d corridors\n", zoneID, len(corridors))

	alert := map[string]interface{}{
		"alert_type": "EVACUATION_REQUIRED",
		"severity":   "CRITICAL",
		"zone_id":    zoneID,
		"message":    "Fire spreading! Use designated evacuation corridors.",
		"corridors":  corridors,
		"timestamp":  getCurrentTimestamp(),
	}

	sendAlert(alert)
	publishToKafka("citizen_alerts", alert)

	fmt.Printf("[Citizen Alert] Evacuation corridor alert sent for zone %s\n", zoneID)
}

func handlePreventionAlert(payload map[string]interface{}) {
	zoneID, _ := payload["zone_id"].(string)
	riskScore, _ := payload["risk_score"].(float64)

	fmt.Printf("[Citizen Alert] Prevention alert for zone %s (risk: %.2f)\n", zoneID, riskScore)

	alert := map[string]interface{}{
		"alert_type": "PREVENTION_ADVISORY",
		"severity":   "MEDIUM",
		"zone_id":    zoneID,
		"risk_score": riskScore,
		"message":    "Fire risk elevated. Avoid outdoor burning. Keep emergency kit ready.",
		"timestamp":  getCurrentTimestamp(),
	}

	sendAlert(alert)
	publishToKafka("citizen_alerts", alert)

	fmt.Printf("[Citizen Alert] Prevention advisory sent for zone %s\n", zoneID)
}

func sendAlert(alert map[string]interface{}) {
	fmt.Printf("[Citizen Alert] Dispatching alert: %s | Severity: %s\n", alert["alert_type"], alert["severity"])

	// Async send notifications
	go sendSMS(alert)
	go sendEmail(alert)
	go sendPushNotification(alert)
}

func sendSMS(alert map[string]interface{}) {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNum := os.Getenv("TWILIO_FROM_NUMBER")
	toNum := os.Getenv("CITIZEN_ALERT_PHONE_NUMBER")

	if accountSid == "" || authToken == "" || fromNum == "" || toNum == "" {
		fmt.Printf("[Citizen Alert] Twilio SMS config incomplete. Logging mock SMS: %s\n", alert["message"])
		return
	}

	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)
	v := url.Values{}
	v.Set("To", toNum)
	v.Set("From", fromNum)
	v.Set("Body", fmt.Sprintf("[%s - %s] %s", alert["alert_type"], alert["severity"], alert["message"]))

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(v.Encode()))
	if err != nil {
		log.Printf("[Citizen Alert] Failed to build Twilio request: %v", err)
		return
	}
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Citizen Alert] Twilio HTTP client error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("[Citizen Alert] Twilio failed: status code %d", resp.StatusCode)
		return
	}
	fmt.Printf("[Citizen Alert] Real SMS dispatched successfully to %s\n", toNum)
}

func sendPushNotification(alert map[string]interface{}) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		fmt.Printf("[Citizen Alert] Discord Push config incomplete. Logging mock Push: %s\n", alert["message"])
		return
	}

	content := fmt.Sprintf("🚨 **[%s - %s]** %s\n*Zone ID: %v*", alert["alert_type"], alert["severity"], alert["message"], alert["zone_id"])
	payload := map[string]string{
		"content": content,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Citizen Alert] Failed to marshal Discord payload: %v", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("[Citizen Alert] Discord HTTP POST error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("[Citizen Alert] Discord Webhook failed: status code %d", resp.StatusCode)
		return
	}
	fmt.Println("[Citizen Alert] Real Push Notification (Discord Webhook) dispatched successfully")
}

func sendEmail(alert map[string]interface{}) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM")
	toEmail := os.Getenv("CITIZEN_ALERT_EMAIL")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || smtpFrom == "" || toEmail == "" {
		fmt.Printf("[Citizen Alert] SMTP config incomplete. Logging mock Email: %s\n", alert["message"])
		return
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	subject := fmt.Sprintf("[%s] Wildfire Emergency Notice", alert["severity"])
	body := fmt.Sprintf("<h2>Blaze-Guard Fire Notification</h2><p><strong>Status:</strong> %s</p><p><strong>Alert:</strong> %s</p><p><strong>Zone:</strong> %v</p><hr><p>Please stay safe and consult safety corridors.</p>", alert["severity"], alert["message"], alert["zone_id"])

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := smtpHost + ":" + smtpPort
	err := smtp.SendMail(addr, auth, smtpFrom, []string{toEmail}, msg)
	if err != nil {
		log.Printf("[Citizen Alert] SMTP Email error: %v", err)
		return
	}
	fmt.Printf("[Citizen Alert] Real Email dispatched successfully to %s\n", toEmail)
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
		log.Printf("[Citizen Alert] Kafka publish error: %v", err)
	}
}

func getCurrentTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
