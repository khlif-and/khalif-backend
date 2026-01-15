-- Up Migration: Create Playlists

-- Create playlists table
CREATE TABLE IF NOT EXISTS playlists (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    author_type VARCHAR(10) NOT NULL,
    author_id INT NOT NULL,
    thumbnail_file TEXT,
    like_count BIGINT DEFAULT 0,
    listening_count BIGINT DEFAULT 0,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_playlists_uuid ON playlists(uuid);
CREATE INDEX IF NOT EXISTS idx_playlists_author ON playlists(author_type, author_id);
CREATE INDEX IF NOT EXISTS idx_playlists_deleted_at ON playlists(deleted_at);

-- Create playlist_audios junction table
CREATE TABLE IF NOT EXISTS playlist_audios (
    id SERIAL PRIMARY KEY,
    playlist_id INT NOT NULL,
    audio_id INT NOT NULL,
    position INT DEFAULT 0,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_playlist_audios_playlist FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE,
    CONSTRAINT fk_playlist_audios_audio FOREIGN KEY (audio_id) REFERENCES audios(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_playlist_audios_playlist ON playlist_audios(playlist_id);
CREATE INDEX IF NOT EXISTS idx_playlist_audios_audio ON playlist_audios(audio_id);

-- Create playlist_likes table
CREATE TABLE IF NOT EXISTS playlist_likes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    playlist_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_playlist_likes_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_playlist_likes_playlist FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE,
    UNIQUE(user_id, playlist_id)
);

CREATE INDEX IF NOT EXISTS idx_playlist_likes_user ON playlist_likes(user_id);
CREATE INDEX IF NOT EXISTS idx_playlist_likes_playlist ON playlist_likes(playlist_id);

-- Stored Procedure: sp_like_playlist
CREATE OR REPLACE FUNCTION sp_like_playlist(
    p_user_id BIGINT,
    p_playlist_id BIGINT
) RETURNS TABLE (
    already_liked BOOLEAN,
    new_count BIGINT
) AS $$
DECLARE
    v_exists BOOLEAN;
    v_new_count BIGINT;
BEGIN
    SELECT EXISTS(
        SELECT 1 FROM playlist_likes
        WHERE user_id = p_user_id AND playlist_id = p_playlist_id
    ) INTO v_exists;

    IF v_exists THEN
        SELECT like_count INTO v_new_count FROM playlists WHERE id = p_playlist_id;
        RETURN QUERY SELECT TRUE, v_new_count;
    ELSE
        INSERT INTO playlist_likes (user_id, playlist_id, created_at)
        VALUES (p_user_id, p_playlist_id, NOW());

        UPDATE playlists
        SET like_count = like_count + 1
        WHERE id = p_playlist_id
        RETURNING like_count INTO v_new_count;

        RETURN QUERY SELECT FALSE, v_new_count;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Stored Procedure: sp_unlike_playlist
CREATE OR REPLACE FUNCTION sp_unlike_playlist(
    p_user_id BIGINT,
    p_playlist_id BIGINT
) RETURNS TABLE (
    was_liked BOOLEAN,
    new_count BIGINT
) AS $$
DECLARE
    v_deleted INT;
    v_new_count BIGINT;
BEGIN
    DELETE FROM playlist_likes
    WHERE user_id = p_user_id AND playlist_id = p_playlist_id;

    GET DIAGNOSTICS v_deleted = ROW_COUNT;

    IF v_deleted > 0 THEN
        UPDATE playlists
        SET like_count = GREATEST(like_count - 1, 0)
        WHERE id = p_playlist_id
        RETURNING like_count INTO v_new_count;

        RETURN QUERY SELECT TRUE, v_new_count;
    ELSE
        SELECT like_count INTO v_new_count FROM playlists WHERE id = p_playlist_id;
        RETURN QUERY SELECT FALSE, v_new_count;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Stored Procedure: sp_record_playlist_listening
CREATE OR REPLACE FUNCTION sp_record_playlist_listening(
    p_playlist_id BIGINT
) RETURNS TABLE (
    new_count BIGINT
) AS $$
DECLARE
    v_new_count BIGINT;
BEGIN
    UPDATE playlists
    SET listening_count = listening_count + 1
    WHERE id = p_playlist_id
    RETURNING listening_count INTO v_new_count;

    RETURN QUERY SELECT v_new_count;
END;
$$ LANGUAGE plpgsql;

-- Stored Procedure: sp_add_audio_to_playlist
CREATE OR REPLACE FUNCTION sp_add_audio_to_playlist(
    p_playlist_id BIGINT,
    p_audio_id BIGINT,
    p_position INT DEFAULT NULL
) RETURNS TABLE (
    already_exists BOOLEAN,
    final_position INT
) AS $$
DECLARE
    v_exists BOOLEAN;
    v_max_position INT;
    v_final_position INT;
BEGIN
    SELECT EXISTS(
        SELECT 1 FROM playlist_audios
        WHERE playlist_id = p_playlist_id AND audio_id = p_audio_id
    ) INTO v_exists;

    IF v_exists THEN
        SELECT position INTO v_final_position 
        FROM playlist_audios 
        WHERE playlist_id = p_playlist_id AND audio_id = p_audio_id;
        
        RETURN QUERY SELECT TRUE, v_final_position;
    ELSE
        IF p_position IS NULL OR p_position = 0 THEN
            SELECT COALESCE(MAX(position), 0) + 1 INTO v_final_position
            FROM playlist_audios
            WHERE playlist_id = p_playlist_id;
        ELSE
            v_final_position := p_position;
        END IF;

        INSERT INTO playlist_audios (playlist_id, audio_id, position, added_at)
        VALUES (p_playlist_id, p_audio_id, v_final_position, NOW());

        RETURN QUERY SELECT FALSE, v_final_position;
    END IF;
END;
$$ LANGUAGE plpgsql;
