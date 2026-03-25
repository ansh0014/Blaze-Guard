# Self Agent Logic

## Purpose
The Self Agent is the continuous improvement and governance layer for all AI agents.  
It monitors production behavior, detects model degradation, selects training data, validates new models, and controls safe rollout.

---

## Scope in This Project
The Self Agent does not replace Detection, Prediction, Logistics, or Citizen Alert agents.  
It supervises and improves them using feedback from real incidents.

---

## Agents It Monitors
- **Detection Agent**: fire detection quality (confidence, false alarms, misses)
- **Prediction Agent**: spread/risk quality (forecast error vs actual fire behavior)
- **Logistics Agent**: route and ETA quality (planned vs actual movement)
- **Citizen Alert Agent**: alert effectiveness (delivery, acknowledgment, escalation)

---

## Input Signals
- Kafka event streams from all agents
- Ground truth and operator feedback
- Real incident outcomes (post-incident records)
- External data changes (weather, satellite feed patterns, seasonal shifts)

---

## Core Responsibilities

### 1) Observability
- Track model confidence trends
- Track drift in features and data distributions
- Track latency, failure rate, and missing data

### 2) Error Analysis
- Identify false positives and false negatives
- Compare predicted spread vs actual spread
- Compare planned ETA vs actual response time
- Measure alert quality (timeliness, relevance)

### 3) Data Selection for Retraining
- Prioritize uncertain samples
- Prioritize high-impact zones
- Prioritize disagreement between model output and observed outcome
- Build curated training datasets

### 4) Candidate Model Evaluation
- Train candidate model versions offline
- Validate on:
  - baseline test set
  - recent real-world incidents
  - edge-case scenarios
- Produce comparison report against current production model

### 5) Safe Deployment Control
- Shadow testing first
- Canary rollout in stages
- Promote only if KPIs improve and safety checks pass
- Auto rollback if quality drops

### 6) Governance and Audit
- Log model version for each decision
- Log confidence and decision metadata
- Keep incident-level traceability for compliance and review

---

## Memory Strategy
- **Short-term memory**: recent incident context for immediate adaptation
- **Long-term memory**: historical incident retrieval for pattern matching and hard-case learning

---

## Safety Rules
- No direct uncontrolled model overwrite in production
- Human approval required for major model updates
- Hard guardrails for high-risk zones and critical alerts
- Rollback path must always be available

---

## High-Level Workflow
1. Consume events and outcomes  
2. Compute quality/drift metrics  
3. Detect degradation or learning opportunity  
4. Trigger retraining pipeline  
5. Evaluate candidate model  
6. Shadow/canary deployment  
7. Promote or rollback  
8. Store full audit trail

---

## Success Criteria
- Lower false alarms and missed detections
- Better spread prediction accuracy
- Better route ETA reliability
- Faster and more relevant citizen alerts
- Stable, auditable, and safe model evolution

---

## Current Team Mapping
- Detection and Prediction models are built by ML team
- Backend integration and orchestration connect agent outputs
- Self Agent coordinates model improvement loop and release safety

---

## Non-Goals
- Not a replacement for emergency command decisions
- Not fully autonomous online learning without validation
- Not a single-source truth for ground reality

---

## Final Note
Self Agent is the control layer that turns separate AI agents into a reliable, continuously improving disaster-response system.