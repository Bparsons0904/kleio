-- Create the styluses table
CREATE TABLE IF NOT EXISTS styluses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  manufacturer TEXT,
  expected_lifespan_hours INTEGER,
  purchase_date TIMESTAMP,
  active BOOLEAN DEFAULT FALSE,
  primary_stylus BOOLEAN DEFAULT FALSE,
  model_number TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(manufacturer, model_number)
);

-- Create the play_history table
CREATE TABLE IF NOT EXISTS play_history (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  release_id INTEGER NOT NULL,
  stylus_id INTEGER,
  played_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (release_id) REFERENCES releases(id),
  FOREIGN KEY (stylus_id) REFERENCES styluses(id)
);

-- Create trigger to update the updated_at timestamp when a stylus is updated
CREATE TRIGGER IF NOT EXISTS styluses_updated_at
AFTER UPDATE ON styluses
FOR EACH ROW
BEGIN
  UPDATE styluses SET updated_at = CURRENT_TIMESTAMP
  WHERE id = OLD.id;
END;

-- Create trigger to update the updated_at timestamp when a play history is updated
CREATE TRIGGER IF NOT EXISTS play_history_updated_at
AFTER UPDATE ON play_history
FOR EACH ROW
BEGIN
  UPDATE play_history SET updated_at = CURRENT_TIMESTAMP
  WHERE id = OLD.id;
END;

-- Ensure only one primary stylus
CREATE TRIGGER IF NOT EXISTS ensure_single_primary
AFTER UPDATE OF primary_stylus ON styluses
FOR EACH ROW WHEN NEW.primary_stylus = 1
BEGIN
  UPDATE styluses SET primary_stylus = 0 WHERE id != NEW.id AND primary_stylus = 1;
END;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_play_history_release_id ON play_history(release_id);
CREATE INDEX IF NOT EXISTS idx_play_history_stylus_id ON play_history(stylus_id);
CREATE INDEX IF NOT EXISTS idx_play_history_played_at ON play_history(played_at);
