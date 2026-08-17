package state

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type A2AMessage struct {
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	EventType string                 `json:"event_type"`
	Payload   map[string]interface{} `json:"payload"`
}

type Store interface {
	AddSuccess(msg interface{})
	AddFailed(msg interface{}, reason string)
}

type InMemoryStore struct {
	mu      sync.Mutex
	success []A2AMessage
	failed  []map[string]any
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		success: make([]A2AMessage, 0, 100),
		failed:  make([]map[string]any, 0, 100),
	}
}

func (s *InMemoryStore) AddSuccess(msg interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m, ok := msg.(A2AMessage); ok {
		s.success = append(s.success, m)
	} else if m, ok := msg.(struct {
		From      string
		To        string
		EventType string
		Payload   map[string]interface{}
	}); ok {
		s.success = append(s.success, A2AMessage{
			From:      m.From,
			To:        m.To,
			EventType: m.EventType,
			Payload:   m.Payload,
		})
	}
}

func (s *InMemoryStore) AddFailed(msg interface{}, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failed = append(s.failed, map[string]any{
		"message": msg,
		"reason":  reason,
	})
}

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(redisURL string) *RedisStore {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("[Redis] Invalid REDIS_URL: %v", err)
	}

	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("[Redis] Warning: Failed to connect to Redis Cloud. Error: %v", err)
	} else {
		log.Println("[Redis] Successfully connected to Redis Cloud")
	}

	return &RedisStore{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisStore) AddSuccess(msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[Redis] Marshal error: %v", err)
		return
	}

	err = r.client.RPush(r.ctx, "blazeguard:a2a:success", data).Err()
	if err != nil {
		log.Printf("[Redis] AddSuccess failed: %v", err)
	}
}

func (r *RedisStore) AddFailed(msg interface{}, reason string) {
	entry := map[string]any{
		"message":   msg,
		"reason":    reason,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("[Redis] Marshal error: %v", err)
		return
	}

	err = r.client.RPush(r.ctx, "blazeguard:a2a:failed", data).Err()
	if err != nil {
		log.Printf("[Redis] AddFailed failed: %v", err)
	}
}

// Factory function to return the correct store based on env
func GetStore() Store {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		return NewRedisStore(redisURL)
	}
	log.Println("[Store] REDIS_URL not set. Using InMemoryStore.")
	return NewInMemoryStore()
}
