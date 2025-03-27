CREATE TABLE IF NOT EXISTS syncs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  sync_start TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  sync_end TIMESTAMP,
  status TEXT DEFAULT 'in_progress', -- 'in_progress', 'complete', 'failed'
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_syncs_sync_start ON syncs(sync_start DESC);
