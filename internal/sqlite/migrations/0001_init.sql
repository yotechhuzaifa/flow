PRAGMA journal_mode = WAL;

CREATE TABLE IF NOT EXISTS flows (
  id         TEXT PRIMARY KEY,
  name       TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  data       TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS runs (
  id         TEXT PRIMARY KEY,
  flow_id    TEXT NOT NULL REFERENCES flows(id) ON DELETE CASCADE,
  status     TEXT NOT NULL,
  started_at TEXT NOT NULL,
  finished_at TEXT
);
