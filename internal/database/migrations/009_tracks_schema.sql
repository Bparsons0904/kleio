CREATE TABLE IF NOT EXISTS tracks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  release_id INTEGER NOT NULL,
  position TEXT NOT NULL,
  title TEXT NOT NULL,
  duration_text TEXT,    
  duration_seconds INTEGER, 
  FOREIGN KEY (release_id) REFERENCES releases(id)
);

