-- Seed data for fire stations in the New Delhi/NCR region
-- Migration: 002_seed_fire_stations.sql

INSERT INTO fire_stations (station_code, station_name, latitude, longitude, is_active, available_trucks, available_ambulances)
VALUES
    ('FS-CP', 'Connaught Place Fire Station', 28.6304, 77.2177, TRUE, 5, 2),
    ('FS-CPY', 'Chanakyapuri Fire Station', 28.5912, 77.1950, TRUE, 4, 2),
    ('FS-LN', 'Laxmi Nagar Fire Station', 28.6360, 77.2789, TRUE, 3, 1),
    ('FS-NP', 'Nehru Place Fire Station', 28.5494, 77.2516, TRUE, 6, 3),
    ('FS-DWK', 'Dwarka Sector 6 Fire Station', 28.5921, 77.0620, TRUE, 4, 2)
ON CONFLICT (station_code) 
DO UPDATE SET 
    latitude = EXCLUDED.latitude, 
    longitude = EXCLUDED.longitude, 
    is_active = EXCLUDED.is_active,
    available_trucks = EXCLUDED.available_trucks,
    available_ambulances = EXCLUDED.available_ambulances,
    updated_at = NOW();
