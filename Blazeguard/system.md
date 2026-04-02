# Blaze-Guard System Architecture

Blaze-Guard is an event-driven wildfire response platform with a microservice architecture.  
It combines ML detection/prediction, geospatial logistics, and citizen alerting.

---

## 1. Architecture Style

- **Data plane:** Kafka topics for asynchronous event flow
- **Control plane:** Orchestrator (HTTP + gRPC) for routing/coordination
- **Service pattern:** A2A agents with independent responsibilities
- **Deployment mode (current):** Local-first with Docker
- **Deployment mode (target):** Cloud-ready (DB/Redis/API keys via env)

---

## 2. Core Components

## 2.1 Data Sources + Ingestion
- NASA FIRMS (event/geo signals)
- Satellite/weather streams
- Ingestion + preprocessing pipeline
- Event builder emits normalized fire events

## 2.2 ML Pipeline
- Detection flow (fire confidence from incoming signals)
- Prediction flow (risk score + spread context)
- ML service (`ml-model/main.py`) integrated as model runtime

## 2.3 Kafka Event Bus
Primary topics:
- `fire_detected`
- `fire_prevention_check`
- `logistics_routes`
- `events_dlq` (dead-letter topic)

Kafka is the primary backbone for inter-service event delivery.

## 2.4 Orchestrator
Responsibilities:
- Agent registry/discovery
- Health monitoring
- Message routing
- gRPC service (`orchestrator/proto/orchestrator.proto`)
- HTTP fallback endpoints

## 2.5 A2A Agents
- **Detection Agent:** receives/normalizes detection events
- **Prediction Agent:** computes spread/risk outputs
- **Logistics Agent:** nearest station + routing (PostGIS + Mapbox)
- **Citizen Alert Agent:** sends prevention/emergency notifications
- **Self Agent:** feedback and continuous improvement signals

## 2.6 API Gateway
- Public REST entrypoint
- Publishes incoming events to Kafka
- Handles client-side integration from frontend

## 2.7 Frontend + Auth
- React dashboard (authority + citizen views)
- Auth backend (Firebase credential-based integration)
- Real-time map/alert UI

---

## 3. Data and Memory Layer

## 3.1 Operational Data
- **PostgreSQL + PostGIS** (current: Docker)
  - Fire stations
  - Geospatial queries
  - Route planning context

## 3.2 Cache/State
- **Redis** (current: Docker)
- **Redis Cloud** (planned for deployment)

## 3.3 AI Agent Memory
- **Oracle AI Database 26ai** for long-term AI memory/vector-aware workflows
  - Used by AI-centric flows (especially prediction/self-evolving logic)

---

## 4. Current Runtime Plan (Implemented Direction)

## 4.1 Local Development (Now)
- Dockerized PostGIS + Redis
- Local Kafka
- Services configured via root `.env`
- Agents consume/publish directly via Kafka

## 4.2 Cloud Deployment (Later)
- DB migration to Neon/Supabase (if selected for relational workloads)
- Redis migration to Redis Cloud
- Secrets managed via deployment environment/secret manager
- Same code paths, env-only switch

---

## 5. Event Contract Standard

Every event must include:
- `event_id`
- `event_version` (current: `v1`)
- `timestamp` (RFC3339)

Domain keys by event type:
- fire location (`latitude`, `longitude`, `zone_id`)
- risk outputs (`risk_score`, `spread_*`)
- logistics outputs (`route`, `deployment`, `safe_zones`)

Reference: `MD/contracts.md`

---

## 6. Reliability & Safety Controls

- Retry/backoff for external API calls
- Timeout-based clients
- Dead-letter handling via `events_dlq`
- Health endpoints per service
- Orchestrator health monitor for agent availability
- Startup-time required env validation

---

## 7. Security and Secrets

- Secrets must not be committed
- `.env.example` is the config template
- Real keys only in local/deployment env files
- Firebase credential files must stay untracked
- API keys in use:
  - Mapbox
  - Weather provider
  - NASA FIRMS
  - Oracle AI DB credentials
  - Redis/DB credentials

---

## 8. End-to-End Flow

1. Ingestion/ML emits fire-related event  
2. API Gateway or pipeline publishes to Kafka  
3. Detection + Prediction agents enrich risk intelligence  
4. Logistics agent computes response route + safe zones  
5. Citizen Alert agent sends alerts (emergency/prevention)  
6. Self agent records outcomes/feedback for improvement  
7. Dashboard visualizes state and alerts in near real-time

---

## 9. Implementation Status Summary

- Architecture foundation is in place
- Event-driven multi-agent flow is active
- gRPC orchestrator scaffold exists
- Docker local infra model is correct
- Remaining work is production hardening (contracts, reliability, cloud ops)