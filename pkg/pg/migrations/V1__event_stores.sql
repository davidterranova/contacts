CREATE SCHEMA IF NOT EXISTS event_store;
SET SCHEMA 'event_store';

CREATE TABLE IF NOT EXISTS events (
  event_id UUID PRIMARY KEY,
  event_type TEXT NOT NULL,
  event_issued_by TEXT NOT NULL,
  event_issued_at TIMESTAMP NOT NULL,

  aggregate_id UUID NOT NULL,
  aggregate_type TEXT NOT NULL,

  data JSONB
);
