package monitor

import (
    "net/http"
    "time"

    "orchestrator/registry"
)

type HealthMonitor struct {
    registry *registry.AgentRegistry
    interval time.Duration
    client   *http.Client
}

func NewHealthMonitor(reg *registry.AgentRegistry, interval time.Duration) *HealthMonitor {
    return &HealthMonitor{
        registry: reg,
        interval: interval,
        client:   &http.Client{Timeout: 3 * time.Second},
    }
}

func (h *HealthMonitor) Start() {
    t := time.NewTicker(h.interval)
    defer t.Stop()

    for range t.C {
        agents := h.registry.List()
        for _, a := range agents {
            resp, err := h.client.Get(a.BaseURL + "/health")
            if err != nil {
                h.registry.SetHealth(a.Name, false, err.Error())
                continue
            }
            _ = resp.Body.Close()
            h.registry.SetHealth(a.Name, resp.StatusCode == http.StatusOK, "")
        }
    }
}