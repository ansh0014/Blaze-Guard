package state

import "sync"

type A2AMessage struct {
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	EventType string                 `json:"event_type"`
	Payload   map[string]interface{} `json:"payload"`
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

	if m, ok := msg.(struct {
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
