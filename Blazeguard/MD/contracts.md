# Event Contracts (Frozen)

## Version
- `event_version`: `v1` (mandatory in all events)

## Required Common Fields
- `event_id` (string)
- `event_version` (string)
- `timestamp` (RFC3339 string)
- `zone_id` (string where applicable)

## Core Events

### FIRE_DETECTED
Required:
- `event_id`, `event_version`, `zone_id`, `latitude`, `longitude`, `timestamp`

### FIRE_SPREAD_PREDICTION
Required:
- `event_id`, `event_version`, `zone_id`, `risk_score`, `timestamp`

### PREPOSITION_RESOURCES
Required:
- `event_id`, `event_version`, `zone_id`, `risk_score`, `timestamp`

### DEPLOYMENT_ROUTE
Required:
- `event_id`, `event_version`, `zone_id`, `timestamp`

### SAFE_ZONES_UPDATE
Required:
- `event_id`, `event_version`, `zone_id`, `safe_zones`, `timestamp`

### PREVENTION_ALERT
Required:
- `event_id`, `event_version`, `zone_id`, `risk_score`, `timestamp`