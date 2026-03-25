package registry

import (
	"os"
	"sync"
	"time"
)

type Agent struct {
	Name        string    `json:"name"`
	BaseURL     string    `json:"base_url"`
	Health      bool      `json:"health"`
	LastChecked time.Time `json:"last_checked"`
	LastError   string    `json:"last_error,omitempty"`
}
type AgentRegistry struct {
	mu     sync.RWMutex
	agents map[string]*Agent
}

func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents: make(map[string]*Agent),
	}
}
func (r *AgentRegistry) Register(name, baseURL string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[name] = &Agent{
		Name:    name,
		BaseURL: baseURL,
		Health:  false,
	}
}
func (r *AgentRegistry) Get(name string) (*Agent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.agents[name]
	return a, ok

}
func (r *AgentRegistry) List() []Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Agent, 0, len(r.agents))
	for _, a := range r.agents {
		out = append(out, *a)
	}
	return out
}
func (r *AgentRegistry) SetHealth(name string, health bool, errMsg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if a, ok := r.agents[name]; ok {
		a.Health = health
		a.LastChecked = time.Now()
		a.LastError = errMsg
	}
}

func (r *AgentRegistry) BootstrapDefaults() {
	r.Register("DETECTION_AGENT", envOr("DETECTION_AGENT_URL", "http://localhost:8001"))
	r.Register("PREDICTION_AGENT", envOr("PREDICTION_AGENT_URL", "http://localhost:8002"))
	r.Register("LOGISTICS_AGENT", envOr("LOGISTICS_AGENT_URL", "http://localhost:8003"))
	r.Register("CITIZEN_ALERT_AGENT", envOr("CITIZEN_ALERT_AGENT_URL", "http://localhost:8004"))
	r.Register("SELF_AGENT", envOr("SELF_AGENT_URL", "http://localhost:8005"))
}

func envOr(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
