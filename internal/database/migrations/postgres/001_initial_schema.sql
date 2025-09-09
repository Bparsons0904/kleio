-- PostgreSQL 18 Initial Schema with UUID v7 Support
-- This migration creates the complete multi-user schema from scratch

-- Enable UUID generation (PostgreSQL 18 has native UUID v7 support)
-- No extensions needed for UUID v7 in PostgreSQL 18+

-- Create users table (multi-user support)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    discogs_username VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create auth_tokens table
CREATE TABLE auth_tokens (
    token VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    discogs_token TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE
);

-- Create shared data tables (releases, artists, labels, etc.)
CREATE TABLE releases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    discogs_id INTEGER UNIQUE NOT NULL,
    title TEXT NOT NULL,
    year INTEGER,
    resource_url TEXT,
    thumb TEXT,
    cover_image TEXT,
    play_duration INTEGER,
    play_duration_estimated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE artists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    discogs_id INTEGER UNIQUE NOT NULL,
    name TEXT NOT NULL,
    resource_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    discogs_id INTEGER UNIQUE NOT NULL,
    name TEXT NOT NULL,
    resource_url TEXT,
    entity_type TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE genres (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE styles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE tracks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    release_id UUID NOT NULL REFERENCES releases(id) ON DELETE CASCADE,
    position VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    duration_text VARCHAR(20),
    duration_seconds INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create many-to-many relationship tables
CREATE TABLE release_artists (
    release_id UUID REFERENCES releases(id) ON DELETE CASCADE,
    artist_id UUID REFERENCES artists(id) ON DELETE CASCADE,
    PRIMARY KEY (release_id, artist_id)
);

CREATE TABLE release_labels (
    release_id UUID REFERENCES releases(id) ON DELETE CASCADE,
    label_id UUID REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (release_id, label_id)
);

CREATE TABLE release_genres (
    release_id UUID REFERENCES releases(id) ON DELETE CASCADE,
    genre_id UUID REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (release_id, genre_id)
);

CREATE TABLE release_styles (
    release_id UUID REFERENCES releases(id) ON DELETE CASCADE,
    style_id UUID REFERENCES styles(id) ON DELETE CASCADE,
    PRIMARY KEY (release_id, style_id)
);

-- Create user-specific tables
CREATE TABLE user_releases (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    release_id UUID REFERENCES releases(id) ON DELETE CASCADE,
    instance_id INTEGER,
    folder_id INTEGER,
    rating INTEGER DEFAULT 0,
    notes TEXT,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, release_id)
);

CREATE TABLE play_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    release_id UUID NOT NULL REFERENCES releases(id) ON DELETE CASCADE,
    stylus_id UUID REFERENCES styluses(id) ON DELETE SET NULL,
    played_at TIMESTAMP WITH TIME ZONE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE styluses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    manufacturer VARCHAR(255),
    expected_lifespan INTEGER DEFAULT 0,
    purchase_date DATE,
    active BOOLEAN DEFAULT FALSE,
    "primary" BOOLEAN DEFAULT FALSE,
    owned BOOLEAN DEFAULT FALSE,
    base_model BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE cleaning_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    release_id UUID NOT NULL REFERENCES releases(id) ON DELETE CASCADE,
    cleaned_at TIMESTAMP WITH TIME ZONE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE folders (
    id INTEGER PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    count INTEGER DEFAULT 0,
    resource_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE syncs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sync_start TIMESTAMP WITH TIME ZONE NOT NULL,
    sync_end TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'in_progress',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_auth_tokens_user_id ON auth_tokens(user_id);
CREATE INDEX idx_releases_discogs_id ON releases(discogs_id);
CREATE INDEX idx_releases_title ON releases(title);
CREATE INDEX idx_artists_discogs_id ON artists(discogs_id);
CREATE INDEX idx_artists_name ON artists(name);
CREATE INDEX idx_labels_discogs_id ON labels(discogs_id);
CREATE INDEX idx_labels_name ON labels(name);
CREATE INDEX idx_tracks_release_id ON tracks(release_id);
CREATE INDEX idx_tracks_position ON tracks(position);
CREATE INDEX idx_user_releases_user_id ON user_releases(user_id);
CREATE INDEX idx_user_releases_release_id ON user_releases(release_id);
CREATE INDEX idx_user_releases_folder_id ON user_releases(folder_id);
CREATE INDEX idx_play_history_user_id ON play_history(user_id);
CREATE INDEX idx_play_history_release_id ON play_history(release_id);
CREATE INDEX idx_play_history_played_at ON play_history(played_at);
CREATE INDEX idx_styluses_user_id ON styluses(user_id);
CREATE INDEX idx_styluses_active ON styluses(active);
CREATE INDEX idx_styluses_primary ON styluses("primary");
CREATE INDEX idx_cleaning_history_user_id ON cleaning_history(user_id);
CREATE INDEX idx_cleaning_history_release_id ON cleaning_history(release_id);
CREATE INDEX idx_cleaning_history_cleaned_at ON cleaning_history(cleaned_at);
CREATE INDEX idx_folders_user_id ON folders(user_id);
CREATE INDEX idx_syncs_user_id ON syncs(user_id);
CREATE INDEX idx_syncs_status ON syncs(status);

-- Create triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_releases_updated_at BEFORE UPDATE ON releases
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_artists_updated_at BEFORE UPDATE ON artists
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_labels_updated_at BEFORE UPDATE ON labels
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tracks_updated_at BEFORE UPDATE ON tracks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_releases_updated_at BEFORE UPDATE ON user_releases
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_play_history_updated_at BEFORE UPDATE ON play_history
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_styluses_updated_at BEFORE UPDATE ON styluses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cleaning_history_updated_at BEFORE UPDATE ON cleaning_history
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_folders_updated_at BEFORE UPDATE ON folders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_syncs_updated_at BEFORE UPDATE ON syncs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();