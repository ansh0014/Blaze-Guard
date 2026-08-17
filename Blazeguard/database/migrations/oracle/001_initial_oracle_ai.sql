-- Oracle Database 23ai Schema Migration: AI Memory & Vector Storage
-- Migration: 001_initial_oracle_ai.sql

-- 1. Table for storing semantic vector representations of wildfire incidents and telemetry
CREATE TABLE wildfire_vector_memory (
    id VARCHAR2(36) DEFAULT SYS_GUID() PRIMARY KEY,
    incident_id VARCHAR2(100) NOT NULL,
    latitude NUMBER(9,6) NOT NULL,
    longitude NUMBER(9,6) NOT NULL,
    severity VARCHAR2(20) NOT NULL,
    weather_summary VARCHAR2(4000),
    event_summary VARCHAR2(4000) NOT NULL,
    event_embedding VECTOR(1536, FLOAT32), -- Native Oracle 23ai Vector datatype
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Native Oracle 23ai HNSW Vector Index for Cosine Distance Search
CREATE VECTOR INDEX idx_wildfire_vectors 
ON wildfire_vector_memory (event_embedding)
ORGANIZATION INVERTED;

-- 2. Table for tracking model self-evolution feedback and performance drift
CREATE TABLE self_evolution_history (
    id VARCHAR2(36) DEFAULT SYS_GUID() PRIMARY KEY,
    agent_name VARCHAR2(100) NOT NULL,
    model_version VARCHAR2(50) NOT NULL,
    accuracy_score NUMBER(5,4) NOT NULL,
    drift_detected NUMBER(1) DEFAULT 0 CHECK (drift_detected IN (0,1)), -- Oracle Boolean Proxy
    retraining_triggered NUMBER(1) DEFAULT 0 CHECK (retraining_triggered IN (0,1)),
    raw_metrics_clob CLOB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Indexing for quick lookups on agent performance histories
CREATE INDEX idx_evolution_agent_ver 
ON self_evolution_history (agent_name, model_version);
