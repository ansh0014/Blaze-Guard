# Logistics AI Agent

## Purpose
Logistics AI Agent receives fire events and prediction outputs, then computes the fastest deployable route from nearest fire stations and prepares response plans for rescue teams and citizen safety flow.

## Inputs
- Kafka topics:
  - `fire_detected`
  - `fire_prevention_check` (optional, for pre-positioning)
- A2A events:
  - `FIRE_DETECTED`
  - `FIRE_SPREAD_PREDICTION`
  - `PREPOSITION_RESOURCES`

## Core Flow
1. Receive event from Kafka or A2A.
2. Normalize event payload to a shared fire-location structure.
3. Query nearest active stations from PostGIS using geospatial radius filter.
4. Apply Haversine pre-filter to reduce candidate stations.
5. Call Mapbox Directions API in parallel for candidates.
6. Select minimum ETA route (traffic-aware profile).
7. Build deployment plan (station, units, ETA, distance).
8. Publish logistics result to Kafka.
9. Notify Citizen Alert agent via A2A.

## Real-Time Shortest Distance Logic
- Source point: fire latitude/longitude from event.
- Candidate points: fire station coordinates from PostGIS.
- Stage 1 filter:
  - `ST_DWithin` in PostGIS for coarse radius filtering.
  - Optional Haversine filter in app for additional quick pruning.
- Stage 2 precise routing:
  - Mapbox Directions API with `driving-traffic`.
  - Evaluate duration and distance for each candidate station.
- Decision:
  - Choose route with minimum duration (ETA first, distance second).

## PostGIS Logic
- Store station location in `GEOGRAPHY(POINT,4326)`.
- Use geospatial index (`GIST`) on location.
- Query nearest stations:
  - `ST_DWithin` for radius constraint
  - `ST_Distance` for ordering by proximity

## Mapbox Usage
- Endpoint: Directions API
- Profile: `mapbox/driving-traffic`
- Required query params:
  - origin and destination coordinates
  - `geometries=geojson`
  - `steps=true`
  - access token from env
- Output used:
  - `duration` for ETA
  - `distance` for distance
  - geometry/steps for route display and dispatch plan

## Concurrency Model
- Goroutine per candidate station route call
- WaitGroup to synchronize completion
- Buffered channel for route results
- Select best successful route after all responses
- Continue even if some route calls fail

## Reliability Controls
- Request timeout per Mapbox call
- Retry with bounded attempts
- API quota control using semaphore/rate limiter
- Graceful failure if all routes fail
- Structured logs for each stage

## Outputs
- Kafka:
  - `logistics_routes`
- A2A to Citizen Alert:
  - `DEPLOYMENT_ROUTE`
  - `SAFE_ZONES_UPDATE` (when spread corridors exist)
  - `PREVENTION_ALERT` (for risk pre-positioning)

## Production Checklist
- [ ] PostGIS station query implemented and used (no hardcoded stations)
- [ ] Haversine pre-filter implemented before Mapbox calls
- [ ] Parallel Mapbox route calculation implemented
- [ ] Traffic-aware profile used (`driving-traffic`)
- [ ] Timeout + retry logic implemented
- [ ] API quota limiter implemented
- [ ] Kafka publish and A2A send integrated
- [ ] Error handling for malformed payloads
- [ ] Health endpoint available
- [ ] End-to-end tested with detection and prediction agents

## Completion Status

### Functional Completion (your definition)
Logistics Agent is **Complete** when:
- it starts successfully,
- consumes detection/prediction events,
- computes and selects a route,
- publishes output to Kafka,
- sends A2A message to Citizen Alert,
- and this flow works end-to-end in your pipeline.

### Production Completion
Logistics Agent is **Production-Ready** only after:
- PostGIS real station query is active,
- traffic-aware Mapbox routing is stable,
- timeout/retry/rate-limit are verified,
- malformed payload handling is tested,
- monitoring + health + load tests are passed.