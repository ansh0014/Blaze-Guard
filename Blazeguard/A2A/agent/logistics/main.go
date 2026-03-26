package main

// one thing i have to check the ml model for the risk score and spread prediction is impplemented in the prediction agent
import (
    "blazeguard/agent/logistics/config"
    database "blazeguard/agent/logistics/db"
    s3 "blazeguard/agent/logistics/server"
    "blazeguard/shared"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "strconv"

    "github.com/segmentio/kafka-go"
)

func main() {
    shared.LoadEnv()
    if err := shared.RequireEnv("KAFKA_BROKER", "EVENT_VERSION", "MAPBOX_API_KEY"); err != nil {
        log.Fatal(err)
    }

    fmt.Println("[Logistics Agent] Starting...")

    config.InitRateLimiter()
    database.InitDB()

    s3.SetMessageHandler(handleA2AMessage)
    go consumeTopic("fire_detected", handleFireEvent)
    go consumeTopic("fire_prevention_check", handlePreventionCheck)

    go s3.StartHTTPServer()

    select {}
}

func consumeTopic(topic string, handler func([]byte)) {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  []string{shared.GetEnv("KAFKA_BROKER", "localhost:9092")},
        Topic:    topic,
        GroupID:  "logistics_service_group_" + topic,
        MinBytes: 1,
        MaxBytes: 10e6,
    })
    defer reader.Close()

    fmt.Printf("[Logistics Agent] Listening to %s\n", topic)
    for {
        msg, err := reader.ReadMessage(context.Background())
        if err != nil {
            log.Printf("[Logistics] Kafka error: %v", err)
            continue
        }
        handler(msg.Value)
    }
}

func handleA2AMessage(eventType string, payload map[string]interface{}) {
    switch eventType {
    case "FIRE_DETECTED":
        processFireLogistics(payload)
    case "FIRE_SPREAD_PREDICTION":
        processFireLogistics(payload)
    case "PREPOSITION_RESOURCES":
        handlePrepositionResources(payload)
    default:
        fmt.Printf("[Logistics] Unknown event type: %s\n", eventType)
    }
}

func handleFireEvent(data []byte) {
    var event map[string]interface{}
    if err := json.Unmarshal(data, &event); err != nil {
        log.Printf("[Logistics] JSON parse error: %v", err)
        return
    }

    processFireLogistics(event)
}

func handlePreventionCheck(data []byte) {
    var event map[string]interface{}
    if err := json.Unmarshal(data, &event); err != nil {
        log.Printf("[Logistics] JSON parse error: %v", err)
        return
    }

    handlePrepositionResources(event)
}

func processFireLogistics(event map[string]interface{}) {
    zoneID, _ := event["zone_id"].(string)
    fmt.Printf("[Logistics] Processing fire event in zone %s\n", zoneID)

    lat, okLat := toFloat64(event["latitude"])
    lon, okLon := toFloat64(event["longitude"])
    if !okLat || !okLon {
        log.Printf("[Logistics] invalid latitude/longitude for zone %s", zoneID)
        return
    }

    fireLocation := map[string]float64{
        "latitude":  lat,
        "longitude": lon,
    }

    fireStations := database.GetNearbyFireStations(fireLocation, 50.0)
    bestRoute, err := calculateBestRouteParallel(fireLocation, fireStations)
    if err != nil {
        fmt.Printf("[Logistics] Error calculating route: %v\n", err)
        return
    }

    deploymentPlan := createDeploymentPlan(bestRoute, event)

    corridors, hasCorridors := event["corridors"].([]interface{})
    if hasCorridors {
        safeZones := calculateSafeZones(fireLocation, corridors)

        publishToKafka("logistics_routes", map[string]interface{}{
            "zone_id":    zoneID,
            "route":      bestRoute,
            "deployment": deploymentPlan,
            "safe_zones": safeZones,
            "timestamp":  event["timestamp"],
        })

        go shared.SendToAgent("CITIZEN_ALERT_AGENT", shared.A2AMessage{
            From:      "logistics_agent",
            To:        "citizen_alert_agent",
            EventType: "SAFE_ZONES_UPDATE",
            Payload: map[string]interface{}{
                "zone_id":           zoneID,
                "safe_zones":        safeZones,
                "evacuation_routes": calculateEvacuationRoutes(fireLocation, safeZones),
                "deployment":        deploymentPlan,
            },
        })

        fmt.Printf("[Logistics] Safe zones calculated: %d zones identified\n", len(safeZones))
    } else {
        publishToKafka("logistics_routes", map[string]interface{}{
            "zone_id":    zoneID,
            "route":      bestRoute,
            "deployment": deploymentPlan,
            "timestamp":  event["timestamp"],
        })

        go shared.SendToAgent("CITIZEN_ALERT_AGENT", shared.A2AMessage{
            From:      "logistics_agent",
            To:        "citizen_alert_agent",
            EventType: "DEPLOYMENT_ROUTE",
            Payload: map[string]interface{}{
                "zone_id":     zoneID,
                "route":       bestRoute,
                "deployment":  deploymentPlan,
                "eta_minutes": bestRoute.Duration / 60,
            },
        })
    }

    fmt.Printf("[Logistics] Route calculated: %s to Fire | ETA: %d min\n",
        bestRoute.StationName, bestRoute.Duration/60)
}

func handlePrepositionResources(payload map[string]interface{}) {
    zoneID, _ := payload["zone_id"].(string)
    riskScore, _ := toFloat64(payload["risk_score"])

    fmt.Printf("[Logistics] Pre-positioning resources for zone %s (risk: %.2f)\n",
        zoneID, riskScore)

    go shared.SendToAgent("CITIZEN_ALERT_AGENT", shared.A2AMessage{
        From:      "logistics_agent",
        To:        "citizen_alert_agent",
        EventType: "PREVENTION_ALERT",
        Payload: map[string]interface{}{
            "zone_id":    zoneID,
            "risk_score": riskScore,
        },
    })
}

func createDeploymentPlan(route database.RouteInfo, event map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "station_id":    route.StationID,
        "station_name":  route.StationName,
        "fire_trucks":   3,
        "ambulances":    2,
        "personnel":     15,
        "equipment":     []string{"hoses", "breathing_apparatus", "axes", "water_tank"},
        "estimated_eta": route.Duration / 60,
        "distance_km":   route.Distance / 1000,
    }
}

func calculateSafeZones(fireLocation map[string]float64, corridors []interface{}) []map[string]interface{} {
    return []map[string]interface{}{
        {
            "zone_id":   "SAFE_001",
            "name":      "Community Center North",
            "latitude":  fireLocation["latitude"] + 0.05,
            "longitude": fireLocation["longitude"] - 0.03,
            "capacity":  500,
        },
        {
            "zone_id":   "SAFE_002",
            "name":      "Stadium West",
            "latitude":  fireLocation["latitude"] - 0.02,
            "longitude": fireLocation["longitude"] - 0.05,
            "capacity":  1000,
        },
    }
}

func calculateEvacuationRoutes(fireLocation map[string]float64, safeZones []map[string]interface{}) []map[string]interface{} {
    routes := []map[string]interface{}{}
    for _, zone := range safeZones {
        routes = append(routes, map[string]interface{}{
            "from":           fireLocation,
            "to":             zone,
            "estimated_time": 15,
        })
    }
    return routes
}

func publishToKafka(topic string, payload map[string]interface{}) {
    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers: []string{shared.GetEnv("KAFKA_BROKER", "localhost:9092")},
        Topic:   topic,
    })
    defer writer.Close()

    data, _ := json.Marshal(payload)
    err := writer.WriteMessages(context.Background(),
        kafka.Message{Value: data},
    )
    if err != nil {
        log.Printf("[Logistics] Kafka publish error: %v", err)
    }
}

func toFloat64(v interface{}) (float64, bool) {
    switch t := v.(type) {
    case float64:
        return t, true
    case float32:
        return float64(t), true
    case int:
        return float64(t), true
    case int64:
        return float64(t), true
    case string:
        f, err := strconv.ParseFloat(t, 64)
        return f, err == nil
    default:
        return 0, false
    }
}