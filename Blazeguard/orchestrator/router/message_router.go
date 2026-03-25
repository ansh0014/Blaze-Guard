package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"orchestrator/registry"
	"orchestrator/state"
)

type A2AMessage struct {
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	EventType string                 `json:"event_type"`
	Payload   map[string]interface{} `json:"payload"`
}

type MessageRouter struct {
	registry *registry.AgentRegistry
	store    *state.InMemoryStore
	client   *http.Client
}

func NewMessageRouter(reg *registry.AgentRegistry, store *state.InMemoryStore) *MessageRouter {
	return &MessageRouter{
		registry: reg,
		store:    store,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (r *MessageRouter) Forward(ctx context.Context, msg A2AMessage) error {
	target, ok := r.registry.Get(msg.To)
	if !ok {
		return fmt.Errorf("target agent not registered: %s", msg.To)
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target.BaseURL+"/receive", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		r.store.AddFailed(msg, err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf("route failed: %s -> %s status=%d", msg.From, msg.To, resp.StatusCode)
		r.store.AddFailed(msg, err.Error())
		return err
	}

	r.store.AddSuccess(msg)
	return nil
}
