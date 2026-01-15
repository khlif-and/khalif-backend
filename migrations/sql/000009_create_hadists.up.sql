-- Up Migration: Create Hadists

CREATE TABLE IF NOT EXISTS hadists (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    nama_hadist VARCHAR(255) NOT NULL,
    perawi_hadist VARCHAR(255),
    nomor_hadist INT DEFAULT 0,
    shahih_status VARCHAR(20) NOT NULL DEFAULT 'shahih',
    kitab_hadist VARCHAR(255),
    
    arabic_hadist TEXT,
    latin_hadist TEXT,
    translate_hadist TEXT,
    description_hadist TEXT,
    audio_hadist TEXT,
    
    category_hadist VARCHAR(100),
    
    like_count BIGINT DEFAULT 0,
    bookmark_count BIGINT DEFAULT 0,
    listening_count BIGINT DEFAULT 0,
    
    source_link TEXT,
    tags VARCHAR(500),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_hadists_uuid ON hadists(uuid);
CREATE INDEX IF NOT EXISTS idx_hadists_category ON hadists(category_hadist);
CREATE INDEX IF NOT EXISTS idx_hadists_kitab ON hadists(kitab_hadist);
CREATE INDEX IF NOT EXISTS idx_hadists_shahih ON hadists(shahih_status);

-- Create Hadist Likes Table
CREATE TABLE IF NOT EXISTS hadist_likes (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    hadist_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_hadist_likes_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_hadist_likes_hadist FOREIGN KEY (hadist_id) REFERENCES hadists(id) ON DELETE CASCADE,
    CONSTRAINT uni_hadist_likes_user_hadist UNIQUE (user_id, hadist_id)
);

CREATE INDEX IF NOT EXISTS idx_hadist_likes_user ON hadist_likes(user_id);
CREATE INDEX IF NOT EXISTS idx_hadist_likes_hadist ON hadist_likes(hadist_id);

-- Create Hadist Bookmarks Table
CREATE TABLE IF NOT EXISTS hadist_bookmarks (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    hadist_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_hadist_bookmarks_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_hadist_bookmarks_hadist FOREIGN KEY (hadist_id) REFERENCES hadists(id) ON DELETE CASCADE,
    CONSTRAINT uni_hadist_bookmarks_user_hadist UNIQUE (user_id, hadist_id)
);

CREATE INDEX IF NOT EXISTS idx_hadist_bookmarks_user ON hadist_bookmarks(user_id);
CREATE INDEX IF NOT EXISTS idx_hadist_bookmarks_hadist ON hadist_bookmarks(hadist_id);
