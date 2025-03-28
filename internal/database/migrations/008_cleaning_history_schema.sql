CREATE TABLE IF NOT EXISTS cleaning_history (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  release_id INTEGER NOT NULL,
  cleaned_at TIMESTAMP NOT NULL,
  notes TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (release_id) REFERENCES releases(id)
);

-- Create trigger to update the updated_at timestamp when a cleaning history is updated
CREATE TRIGGER IF NOT EXISTS cleaning_history_updated_at
AFTER UPDATE ON cleaning_history
FOR EACH ROW
BEGIN
  UPDATE cleaning_history SET updated_at = CURRENT_TIMESTAMP
  WHERE id = OLD.id;
END;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_cleaning_history_release_id ON cleaning_history(release_id);
CREATE INDEX IF NOT EXISTS idx_cleaning_history_cleaned_at ON cleaning_history(cleaned_at);
