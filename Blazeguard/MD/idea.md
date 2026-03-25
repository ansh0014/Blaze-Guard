# BlazeGuard IDEA (E2E Architecture + gRPC Plan)

## 1) Project Goal
Build an AI-driven wildfire response platform that detects fire, predicts spread, computes logistics response, alerts citizens, and continuously improves models through a self-evolving feedback loop.

---

## 2) Team Mapping
- **Mahi**: ML
- **Abhinandan**: ML + GenAI
- **Rishank**: Frontend
- **Anshul**: Backend, Agents, API Gateway, Kafka, Integration, Deployment

---

## 3) Core AI Agents
1. **Detection Agent**  
   Consumes detection events, validates fire incidents, publishes actionable events.
2. **Prediction Agent**  
   Computes risk score + spread/corridors.
3. **Logistics Agent**  
   Finds nearest stations, shortest ETA route (Mapbox), deployment plan, safe zones.
4. **Citizen Alert Agent**  
   Sends incident/risk/evacuation alerts (SMS/Push/Email based on available provider).
5. **Self-Evolving Agent**  
   Monitors model quality/drift, feedback loops, and controlled improvement signals.

---

## 4) Current Communication Model
- **Event backbone**: Kafka
- **Agent-to-Agent protocol**: A2A HTTP (`/receive`)
- **Frontend/API**: API Gateway (HTTP/REST)
- **Auth**: Backend/Auth (Firebase + Postgres)

This model is already valid for MVP.

---

## 5) Where gRPC Fits
Use **hybrid architecture** (recommended):

- **Keep REST**: Frontend ↔ API Gateway
- **Keep Kafka**: async event pipeline
- **Add gRPC**: internal service calls for low-latency orchestration and typed contracts

### gRPC use-cases
- Orchestrator ↔ all agents (control-plane + health + command)
- API Gateway ↔ Orchestrator (optional optimization)
- Internal synchronous calls that need strict schema + speed

### Do not replace
- Browser-facing APIs (keep HTTP/REST)
- Kafka event fanout (still best for decoupled async workflow)

---

## 6) End-to-End Process (Main Incident Flow)

## Step A: Data Ingestion + ML
1. External source (NASA FIRMS / image stream) ingested.
2. ML pipeline creates:
   - detection output (fire confidence, location, timestamp)
   - prediction output (risk score, spread corridors)

## Step B: Event Distribution
3. ML outputs are published to Kafka topics (e.g., `fire_detected`, `fire_prevention_check`).
4. Detection/Prediction agents process and emit normalized A2A/Kafka events.

## Step C: Logistics Decisioning
5. Logistics agent receives fire/prediction events.
6. Logistics agent:
   - fetches nearby stations from PostGIS
   - optional Haversine pre-filter
   - parallel Mapbox route calls (traffic-aware profile)
   - picks minimum ETA route
   - builds deployment plan + safe-zone payload

## Step D: Citizen Communication
7. Logistics + Prediction trigger Citizen Alert agent.
8. Citizen Alert sends appropriate alert type:
   - critical emergency
   - evacuation corridor/safe zone
   - prevention advisory
9. Alerts also published to Kafka for traceability and UI stream.

## Step E: Self-Evolving Loop
10. Self agent consumes outcomes/feedback from all agents.
11. Tracks quality/drift/confidence/error patterns.
12. Emits model feedback and improvement signals for ML retraining pipeline.

---

## 7) Prevention Flow (Non-emergency)
1. Prediction emits high risk (`fire_prevention_check`).
2. Logistics pre-positions resources.
3. Citizen Alert sends prevention advisory.
4. Self agent logs outcome effectiveness for future tuning.

---

## 8) Orchestrator (“Main Brain”) Role
Your repo already has an `orchestrator/` module. It should be the control-plane:
- Agent registry/discovery
- Health checks
- Workflow routing policy
- Retry/fallback decisions
- Future gRPC control channel

Kafka remains the data-plane for events.

---

## 9) Suggested gRPC Interfaces (High Level)
- `HealthService`: health/readiness checks
- `CommandService`: trigger/ack agent actions
- `FeedbackService`: model quality and drift reports
- `RoutingService` (optional): sync route request/response for critical paths

All protobuf schemas should match existing A2A event payload fields to avoid duplication drift.

---

## 10) Data + Infra
- **PostGIS**: fire stations, geospatial queries
- **Kafka**: event stream + decoupling
- **Mapbox**: shortest/fastest path with traffic
- **pgvector (planned/used)**: long-term memory for self-evolving intelligence
- **LangChain memory**: short-term contextual reasoning

---

## 11) Reliability Requirements
- Topic consumer groups per agent
- Request timeout + retry + rate-limit (Mapbox/API calls)
- Dead-letter strategy for malformed events
- Idempotent event handling
- Structured logs + correlation IDs
- Health endpoints on all services

---

## 12) Functional E2E Definition (Your Rule)
Project is functionally complete when:
- all 5 agents run,
- detection/prediction events trigger downstream flow,
- logistics route and deployment are generated,
- citizen alerts are issued,
- self agent receives feedback,
- complete incident path is visible end-to-end.

---

## 13) Production-Ready Definition
Production-ready only after:
- strict schema contracts (proto + versioning),
- monitored SLOs (latency, error rate, alert delay),
- canary/rollback for model and service updates,
- security hardening (authN/authZ, secrets, TLS),
- load and chaos validation.

---

## 14) Final Architecture Summary
**Best-fit architecture for your project:**
- REST for client APIs
- Kafka for async event pipeline
- A2A (current) for immediate agent interoperability
- gRPC for internal control-plane and typed low-latency service communication via orchestrator

This gives practical MVP speed now and scalable microservice maturity next.