# Work Left (Current Status)

## Overall Progress
- **Estimated done:** 70%
- **Estimated left:** 30%

---

## 1) Critical Left (Must finish before deployment)

### A. Data + Infra
- [ ] Move Logistics DB to cloud-ready config (PostGIS connection via env)
- [ ] Add migration + seed for fire stations
- [ ] Integrate **Redis Cloud** in orchestrator state (TLS/auth)
- [ ] Finalize Kafka topics + DLQ topics

### B. Secrets & Env
- [ ] Centralize environment variables across all services
- [ ] Remove hardcoded values and local-only assumptions
- [ ] Add `.env.example` for each service
- [ ] Validate startup fails clearly when key env is missing

### C. Contracts
- [ ] Freeze Kafka/A2A payload schema (shared contract doc)
- [ ] Add payload validation in all agents
- [ ] Add versioning for event schema (`event_version`)

### D. Reliability
- [ ] Retry/backoff + timeout policy in all external calls
- [ ] Dead-letter handling for malformed/unprocessable events
- [ ] Idempotency handling (duplicate events)

---

## 2) High Priority Left (Should finish for production-ready MVP)

### A. Orchestrator + gRPC
- [ ] Complete gRPC integration tests
- [ ] Keep HTTP fallback stable
- [ ] Health/readiness checks per agent from orchestrator

### B. Logistics
- [ ] Verify real shortest-time routing with Mapbox traffic profile
- [ ] Confirm PostGIS nearest-station query works with real data
- [ ] Add failure fallback when Mapbox unavailable

### C. Citizen Alert
- [ ] Replace mock alert senders with real provider integration (SMS/Push/Email)
- [ ] Add alert deduplication and escalation rules

### D. Self Agent
- [ ] Add drift/quality metrics persistence (Redis/DB)
- [ ] Add threshold rules for feedback events
- [ ] Add basic evaluation report publishing

---

## 3) Medium Priority Left

- [ ] End-to-end tests: Detection → Prediction → Logistics → Citizen Alert → Self
- [ ] Observability: structured logs, request IDs, metrics dashboard
- [ ] Linux-safe path and import cleanup (case sensitivity)
- [ ] Docker Compose full stack run for one-command local startup

---

## 4) Final Hardening

- [ ] Security review (CORS policy, auth checks, secret handling)
- [ ] Load test Kafka consumers and route generation
- [ ] Deployment manifests (staging/prod)
- [ ] Rollback plan + runbook documentation

---

## Estimated Remaining Effort

- **Core deployment-ready MVP:** 10–14 days  
- **Production hardening:** +10–15 days

---

## Exit Criteria (Project considered complete)
- [ ] All 5 agents running in integrated pipeline
- [ ] API Gateway + Orchestrator routing stable
- [ ] Cloud DB + Redis Cloud connected
- [ ] Real Mapbox + PostGIS routing verified
- [ ] E2E flow validated with test incidents
- [ ] Basic monitoring + error handling in place 