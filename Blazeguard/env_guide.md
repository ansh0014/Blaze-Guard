# Blaze-Guard Environment Setup Guide

This guide details all the environment variables required to run the **Blaze-Guard** system, including the Go A2A microservices, the API Gateway, and the FastAPI ML Model server.

The configuration file is located at `d:\code-2-main\Blaze-Guard\Blazeguard\.env`.

---

## 1. General Config

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `APP_ENV` | `development` | Environment mode (`development` or `production`). |
| `EVENT_VERSION` | `v1` | Version for Kafka event message schemas. |

---

## 2. Kafka Configuration

Used by all agents and the API gateway to publish and consume event streams.

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `KAFKA_BROKER` | `localhost:9092` | Connection string of the Kafka broker. |
| `KAFKA_DLQ_TOPIC` | `events_dlq` | Dead Letter Queue topic name for failed events. |

---

## 3. Microservice Networking & Service Discovery

These ports and URLs map how the Gateway, Orchestrator, and individual agents locate and communicate with each other.

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `PORT` | `8080` | Port of the API Gateway. |
| `ORCHESTRATOR_PORT` | `8090` | HTTP Port of the Orchestrator. |
| `ORCHESTRATOR_GRPC_PORT` | `9090` | gRPC Port of the Orchestrator. |
| `FRONTEND_URL` | `http://localhost:5173` | Address of the Vite React Frontend. |
| `VITE_API_BASE_URL` | `http://localhost:8080` | Frontend environment variable to target the gateway. |
| `VITE_AUTH_BASE_URL` | `http://localhost:2000` | Frontend environment variable to target auth endpoints. |
| `VITE_ORCHESTRATOR_URL` | `http://localhost:8090`| Frontend environment variable to target the orchestrator. |
| `DETECTION_AGENT_URL` | `http://localhost:8001` | Address of the Detection agent. |
| `PREDICTION_AGENT_URL` | `http://localhost:8002` | Address of the Prediction agent. |
| `LOGISTICS_AGENT_URL` | `http://localhost:8003` | Address of the Logistics agent. |
| `CITIZEN_ALERT_AGENT_URL`| `http://localhost:8004` | Address of the Citizen Alert agent. |
| `SELF_AGENT_URL` | `http://localhost:8005` | Address of the Self-Evolving agent. |

---

## 4. Third-Party API Keys

These must be replaced with your active API keys to run routing, weather calculations, and satellite downloads.

| Variable | Default Value | Description / Instructions |
| :--- | :--- | :--- |
| `MAPBOX_API_KEY` | `REPLACE_ME` | Go to [Mapbox](https://mapbox.com/) to get a token. Required for route duration calculations in the Logistics agent. |
| `MAPBOX_BASE_URL` | `https://api.mapbox.com` | Base URL for Mapbox directions API. |
| `WEATHER_API_KEY` | `REPLACE_ME` | Go to [OpenWeatherMap](https://openweathermap.org/) to get an API key. Used to fetch wind speed, temperature, and humidity. |
| `WEATHER_API_BASE_URL` | `https://api.openweathermap.org`| Base URL for OpenWeatherMap requests. |
| `NASA_FIRMS_API_KEY` | `REPLACE_ME` | Map transaction token for NASA FIRMS real-time satellite telemetry ingestion. |

---

## 5. Relational & Cache Database Config

PostGIS is used by the Logistics agent to perform spatial searches of fire stations; Redis stores active state.

| Variable | Default / Example Value | Description |
| :--- | :--- | :--- |
| `DB_HOST` | `localhost` | PostgreSQL host. |
| `DB_PORT` | `5432` | PostgreSQL port. |
| `DB_USER` | `postgres` | PostgreSQL username. |
| `DB_PASSWORD` | `REPLACE_ME` | PostgreSQL password. |
| `DB_NAME` | `blazeguard` | PostgreSQL database name (needs PostGIS extension enabled). |
| `DB_SSLMODE` | `disable` | SSL mode (`disable`, `require`, `verify-full`). |
| `DATABASE_URL` | `postgres://postgres:pass@localhost:5432/blazeguard?sslmode=disable` | Combined connection string. |
| `REDIS_URL` | `rediss://default:REPLACE_ME@your-redis-host:6380` | Combined Redis connection URL (SSL). |
| `REDIS_PASSWORD` | `REPLACE_ME` | Password for Redis Cloud instances. |

## 6. Oracle Database 23ai AI Memory Config

Used by the Self-Evolving Agent and long-term storage to handle semantic event summaries, agent feedback loops, and vector indexing.

| Variable | Default / Example Value | Description |
| :--- | :--- | :--- |
| `ORACLE_AI_DB_HOST` | `localhost` | Oracle Database listener hostname. |
| `ORACLE_AI_DB_PORT` | `1521` | Oracle Database port (standard: 1521). |
| `ORACLE_AI_DB_SERVICE` | `FREEPDB1` | Pluggable Database (PDB) service name. |
| `ORACLE_AI_DB_USER` | `blazeguard` | Database schema user. |
| `ORACLE_AI_DB_PASSWORD` | `REPLACE_ME` | Password for the schema user. |
| `ORACLE_AI_DB_DSN` | `oracle://blazeguard:pass@localhost:1521/FREEPDB1` | DSN connection string (for oracle drivers/clients). |

---

## 7. Machine Learning Model Configuration

Configures the FastAPI Python ML service hosting YOLO and severity prediction.

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `ML_MODEL_HOST` | `0.0.0.0` | IP to bind the FastAPI server. |
| `ML_MODEL_PORT` | `9000` | Port of the FastAPI server. |
| `ML_MODEL_URL` | `http://localhost:9000` | URL where the Go services contact the ML models. |
| `ML_MODEL_NAME` | `blazeguard-fire-model` | Human-readable model identifier. |
| `MODEL_PATH` | `./models/fire_model.pt`| Path to the YOLOv8 weights file on disk. |
| `ML_MODEL_API_KEY` | `REPLACE_ME` | API auth token for ML inference calls (optional/development). |

---

## 8. Reliability & Resiliency Constants

Constants tuning timeout and delay thresholds for microservice interactions.

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `REQUEST_TIMEOUT_SECONDS`| `10` | Maximum time to wait for agent HTTP A2A response before timing out. |
| `RETRY_MAX_ATTEMPTS` | `3` | Maximum number of route retrieval retries. |
| `RETRY_BASE_DELAY_MS` | `300` | Exponential backoff baseline multiplier. |

---

## 9. Notification Services Configuration

Enables the **Citizen Alert Agent** to dispatch real-time emergency, risk, and deployment alerts over cellular network, email, and discord integrations.

| Variable | Default / Example Value | Description |
| :--- | :--- | :--- |
| `TWILIO_ACCOUNT_SID` | `REPLACE_ME` | Twilio Account SID for SMS integration. |
| `TWILIO_AUTH_TOKEN` | `REPLACE_ME` | Twilio Auth Token for SMS integration. |
| `TWILIO_FROM_NUMBER` | `REPLACE_ME` | Approved Twilio SMS source phone number. |
| `CITIZEN_ALERT_PHONE_NUMBER` | `REPLACE_ME` | Targeted phone number for emergency SMS alerts. |
| `SMTP_HOST` | `smtp.gmail.com` | SMTP host for dispatching emails. |
| `SMTP_PORT` | `587` | SMTP port (standard: 587/TLS or 465/SSL). |
| `SMTP_USER` | `REPLACE_ME` | SMTP username / email address. |
| `SMTP_PASSWORD` | `REPLACE_ME` | SMTP password / app-specific password. |
| `SMTP_FROM` | `REPLACE_ME` | Sender email address. |
| `CITIZEN_ALERT_EMAIL` | `REPLACE_ME` | Targeted email address for emergency alerts. |
| `DISCORD_WEBHOOK_URL` | `REPLACE_ME` | Discord webhook endpoint to receive push notifications. |

