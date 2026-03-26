package shared

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

var EventSchemas = map[string][]string{
	"FIRE_DETECTED":          {"event_id", "event_version", "zone_id", "latitude", "longitude", "timestamp"},
	"FIRE_SPREAD_PREDICTION": {"event_id", "event_version", "zone_id", "risk_score", "timestamp"},
	"PREPOSITION_RESOURCES":  {"event_id", "event_version", "zone_id", "risk_score", "timestamp"},
	"DEPLOYMENT_ROUTE":       {"event_id", "event_version", "zone_id", "timestamp"},
	"SAFE_ZONES_UPDATE":      {"event_id", "event_version", "zone_id", "safe_zones", "timestamp"},
	"PREVENTION_ALERT":       {"event_id", "event_version", "zone_id", "risk_score", "timestamp"},
}

func EnsureCommonFields(eventType string, payload map[string]interface{}) map[string]interface{} {
	if payload == nil {
		payload = map[string]interface{}{}
	}
	if _, ok := payload["timestamp"]; !ok {
		payload["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}
	if _, ok := payload["event_version"]; !ok {
		payload["event_version"] = GetEnv("EVENT_VERSION", "v1")
	}
	if _, ok := payload["event_id"]; !ok {
		payload["event_id"] = makeEventID(eventType, payload)
	}
	return payload
}

func ValidateEvent(eventType string, payload map[string]interface{}) error {
	required, ok := EventSchemas[eventType]
	if !ok {
		return nil
	}
	for _, k := range required {
		if _, exists := payload[k]; !exists {
			return fmt.Errorf("invalid payload: missing key '%s' for event '%s'", k, eventType)
		}
	}
	return nil
}

func makeEventID(eventType string, payload map[string]interface{}) string {
	raw := fmt.Sprintf("%s|%v|%v|%v", eventType, payload["zone_id"], payload["latitude"], payload["timestamp"])
	sum := sha1.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])
}
