-- Down Migration: Create Playlists

DROP FUNCTION IF EXISTS sp_add_audio_to_playlist;
DROP FUNCTION IF EXISTS sp_record_playlist_listening;
DROP FUNCTION IF EXISTS sp_unlike_playlist;
DROP FUNCTION IF EXISTS sp_like_playlist;
DROP INDEX IF EXISTS idx_playlist_likes_playlist;
DROP INDEX IF EXISTS idx_playlist_likes_user;
DROP TABLE IF EXISTS playlist_likes;
DROP INDEX IF EXISTS idx_playlist_audios_audio;
DROP INDEX IF EXISTS idx_playlist_audios_playlist;
DROP TABLE IF EXISTS playlist_audios;
DROP INDEX IF EXISTS idx_playlists_deleted_at;
DROP INDEX IF EXISTS idx_playlists_author;
DROP INDEX IF EXISTS idx_playlists_uuid;
DROP TABLE IF EXISTS playlists;
