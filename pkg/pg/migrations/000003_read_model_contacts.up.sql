-- SET SCHEMA read_model;

CREATE TABLE IF NOT EXISTS read_model.contacts (
  id UUID PRIMARY KEY,
  aggregate_version INT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP,
  created_by TEXT NOT NULL,
  first_name TEXT,
  last_name TEXT,
  email TEXT,
  phone TEXT
);

CREATE INDEX IF NOT EXISTS created_byidx ON read_model.contacts (created_by);
