-- Откат миграции 001_initial_schema

DROP TRIGGER IF EXISTS update_document_templates_updated_at ON document_templates;
DROP TRIGGER IF EXISTS update_blocking_rules_updated_at ON blocking_rules;
DROP TRIGGER IF EXISTS update_tickets_updated_at ON tickets;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_seats_updated_at ON seats;
DROP TRIGGER IF EXISTS update_trips_updated_at ON trips;
DROP TRIGGER IF EXISTS update_schedules_updated_at ON schedules;
DROP TRIGGER IF EXISTS update_routes_updated_at ON routes;
DROP TRIGGER IF EXISTS update_drivers_updated_at ON drivers;
DROP TRIGGER IF EXISTS update_buses_updated_at ON buses;
DROP TRIGGER IF EXISTS update_stations_updated_at ON stations;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS announcements;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS fiscal_receipts;
DROP TABLE IF EXISTS boarding_events;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS document_templates;
DROP TABLE IF EXISTS blocking_rules;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS seats;
DROP TABLE IF EXISTS trips;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS drivers;
DROP TABLE IF EXISTS buses;
DROP TABLE IF EXISTS stations;

DROP EXTENSION IF EXISTS "pg_trgm";
DROP EXTENSION IF EXISTS "uuid-ossp";
