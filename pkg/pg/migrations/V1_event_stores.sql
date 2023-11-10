CREATE SCHEMA IF NOT EXISTS event_store;
SET SCHEMA 'event_store';

CREATE TABLE IF NOT EXISTS events (
  id UUID PRIMARY KEY,
  aggregate_id UUID NOT NULL,
  aggregate_type TEXT NOT NULL,
  event_type TEXT NOT NULL,

  created_at TIMESTAMP NOT NULL,
  created_by TEXT NOT NULL,

  data JSONB
);
