-- Create the releases table to store Discogs release information
CREATE TABLE IF NOT EXISTS releases (
  id INTEGER PRIMARY KEY, -- Release ID from Discogs
  instance_id INTEGER NOT NULL,
  folder_id INTEGER NOT NULL,
  rating INTEGER NOT NULL DEFAULT 0,
  title TEXT NOT NULL,
  year INTEGER,
  resource_url TEXT,
  thumb TEXT, -- Thumbnail image URL
  cover_image TEXT, -- Cover image URL
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_synced TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (folder_id) REFERENCES folders(id)
);

-- Create the labels table to store label information
CREATE TABLE IF NOT EXISTS labels (
  id INTEGER PRIMARY KEY, -- Label ID from Discogs
  name TEXT NOT NULL,
  resource_url TEXT,
  entity_type TEXT
);

-- Create the release_labels junction table
CREATE TABLE IF NOT EXISTS release_labels (
  release_id INTEGER NOT NULL,
  label_id INTEGER NOT NULL,
  catno TEXT, -- Catalog number
  PRIMARY KEY (release_id, label_id, catno),
  FOREIGN KEY (release_id) REFERENCES releases(id),
  FOREIGN KEY (label_id) REFERENCES labels(id)
);

-- Create the artists table
CREATE TABLE IF NOT EXISTS artists (
  id INTEGER PRIMARY KEY, -- Artist ID from Discogs
  name TEXT NOT NULL,
  resource_url TEXT
);

-- Create the release_artists junction table
CREATE TABLE IF NOT EXISTS release_artists (
  release_id INTEGER NOT NULL,
  artist_id INTEGER NOT NULL,
  join_relation TEXT, -- The "join" field in the API
  anv TEXT, -- Artist name variation
  tracks TEXT, -- Which tracks this artist appears on
  role TEXT, -- Artist role (e.g., "Producer")
  PRIMARY KEY (release_id, artist_id, role),
  FOREIGN KEY (release_id) REFERENCES releases(id),
  FOREIGN KEY (artist_id) REFERENCES artists(id)
);

-- Create the formats table
CREATE TABLE IF NOT EXISTS formats (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  release_id INTEGER NOT NULL,
  name TEXT NOT NULL, -- Format name (e.g., "CDr", "Vinyl")
  qty INTEGER NOT NULL DEFAULT 1, -- Quantity of this format
  FOREIGN KEY (release_id) REFERENCES releases(id)
);

-- Create format_descriptions table for the descriptions array
CREATE TABLE IF NOT EXISTS format_descriptions (
  format_id INTEGER NOT NULL,
  description TEXT NOT NULL,
  PRIMARY KEY (format_id, description),
  FOREIGN KEY (format_id) REFERENCES formats(id)
);

-- Create genres table
CREATE TABLE IF NOT EXISTS genres (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT UNIQUE NOT NULL
);

-- Create the release_genres junction table
CREATE TABLE IF NOT EXISTS release_genres (
  release_id INTEGER NOT NULL,
  genre_id INTEGER NOT NULL,
  PRIMARY KEY (release_id, genre_id),
  FOREIGN KEY (release_id) REFERENCES releases(id),
  FOREIGN KEY (genre_id) REFERENCES genres(id)
);

-- Create styles table
CREATE TABLE IF NOT EXISTS styles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT UNIQUE NOT NULL
);

-- Create the release_styles junction table
CREATE TABLE IF NOT EXISTS release_styles (
  release_id INTEGER NOT NULL,
  style_id INTEGER NOT NULL,
  PRIMARY KEY (release_id, style_id),
  FOREIGN KEY (release_id) REFERENCES releases(id),
  FOREIGN KEY (style_id) REFERENCES styles(id)
);

-- Create notes table
CREATE TABLE IF NOT EXISTS release_notes (
  release_id INTEGER NOT NULL,
  field_id INTEGER NOT NULL,
  value TEXT NOT NULL,
  PRIMARY KEY (release_id, field_id),
  FOREIGN KEY (release_id) REFERENCES releases(id)
);

-- Create trigger to update the updated_at timestamp on releases table
CREATE TRIGGER IF NOT EXISTS releases_updated_at
AFTER UPDATE ON releases
FOR EACH ROW
BEGIN
  UPDATE releases SET updated_at = CURRENT_TIMESTAMP
  WHERE id = OLD.id;
END;
