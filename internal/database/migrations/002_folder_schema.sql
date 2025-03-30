CREATE TABLE IF NOT EXISTS folders (
  id INTEGER PRIMARY KEY, -- Using the Discogs folder ID as our primary key
  name TEXT NOT NULL,
  count INTEGER NOT NULL DEFAULT 0,
  resource_url TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger to update the updated_at timestamp when a row is updated
CREATE TRIGGER IF NOT EXISTS folders_updated_at
AFTER UPDATE ON folders
FOR EACH ROW
BEGIN
  UPDATE folders SET updated_at = CURRENT_TIMESTAMP
  WHERE id = OLD.id;
END;
