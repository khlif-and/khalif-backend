-- Down Migration: Create Doas

DROP INDEX IF EXISTS idx_user_doa_bookmark;
DROP TABLE IF EXISTS doa_bookmarks;
DROP INDEX IF EXISTS idx_user_doa_like;
DROP TABLE IF EXISTS doa_likes;
DROP INDEX IF EXISTS idx_doas_deleted_at;
DROP INDEX IF EXISTS idx_doas_hadist_id;
DROP INDEX IF EXISTS idx_doas_category;
DROP INDEX IF EXISTS idx_doas_uuid;
DROP TABLE IF EXISTS doas;
