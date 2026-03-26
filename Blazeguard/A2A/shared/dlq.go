package shared

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type DLQMessage struct {
	SourceService string                 `json:"source_service"`
	SourceTopic   string                 `json:"source_topic"`
	EventType     string                 `json:"event_type"`
	Reason        string                 `json:"reason"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	Raw           string                 `json:"raw,omitempty"`
	Timestamp     string                 `json:"timestamp"`
}

func PublishDLQ(sourceService, sourceTopic, eventType, reason string, payload map[string]interface{}, raw []byte) {
	broker := GetEnv("KAFKA_BROKER", "localhost:9092")
	dlqTopic := GetEnv("KAFKA_DLQ_TOPIC", "events_dlq")

	msg := DLQMessage{
		SourceService: sourceService,
		SourceTopic:   sourceTopic,
		EventType:     eventType,
		Reason:        reason,
		Payload:       payload,
		Raw:           string(raw),
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   dlqTopic,
	})
	defer w.Close()

	if err := w.WriteMessages(context.Background(), kafka.Message{Value: data}); err != nil {
		log.Printf("[DLQ] publish failed: %v", err)
	}
}
