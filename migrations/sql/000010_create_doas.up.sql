-- Up Migration: Create Doas

CREATE TABLE IF NOT EXISTS doas (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    judul_doa VARCHAR(255) NOT NULL,
    arabic_doa TEXT,
    latin_doa TEXT,
    translate_doa TEXT,
    description_doa TEXT,
    audio_doa TEXT,
    category_doa VARCHAR(100),
    source_link VARCHAR(255),
    tags TEXT,
    like_count INTEGER DEFAULT 0,
    bookmark_count INTEGER DEFAULT 0,
    listening_count INTEGER DEFAULT 0,
    
    hadist_id INTEGER REFERENCES hadists(id) ON DELETE SET NULL, 

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_doas_uuid ON doas(uuid);
CREATE INDEX IF NOT EXISTS idx_doas_category ON doas(category_doa);
CREATE INDEX IF NOT EXISTS idx_doas_hadist_id ON doas(hadist_id);
CREATE INDEX IF NOT EXISTS idx_doas_deleted_at ON doas(deleted_at);

-- Create doa_likes table
CREATE TABLE IF NOT EXISTS doa_likes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    doa_id INTEGER NOT NULL REFERENCES doas(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, doa_id)
);

CREATE INDEX IF NOT EXISTS idx_user_doa_like ON doa_likes(user_id, doa_id);

-- Create doa_bookmarks table
CREATE TABLE IF NOT EXISTS doa_bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    doa_id INTEGER NOT NULL REFERENCES doas(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, doa_id)
);

CREATE INDEX IF NOT EXISTS idx_user_doa_bookmark ON doa_bookmarks(user_id, doa_id);
