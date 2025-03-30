-- Index for faster folder-based queries
CREATE INDEX IF NOT EXISTS idx_releases_folder_id ON releases(folder_id);

-- Index for artist lookups
CREATE INDEX IF NOT EXISTS idx_release_artists_artist_id ON release_artists(artist_id);
CREATE INDEX IF NOT EXISTS idx_release_artists_release_id ON release_artists(release_id);

-- Index for label lookups
CREATE INDEX IF NOT EXISTS idx_release_labels_label_id ON release_labels(label_id);
CREATE INDEX IF NOT EXISTS idx_release_labels_release_id ON release_labels(release_id);

-- Index for genre lookups
CREATE INDEX IF NOT EXISTS idx_release_genres_genre_id ON release_genres(genre_id);
CREATE INDEX IF NOT EXISTS idx_release_genres_release_id ON release_genres(release_id);

-- Index for style lookups
CREATE INDEX IF NOT EXISTS idx_release_styles_style_id ON release_styles(style_id);
CREATE INDEX IF NOT EXISTS idx_release_styles_release_id ON release_styles(release_id);

-- Index for faster title searches
CREATE INDEX IF NOT EXISTS idx_releases_title ON releases(title);

-- Index for year-based queries
CREATE INDEX IF NOT EXISTS idx_releases_year ON releases(year);
