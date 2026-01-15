-- Down Migration: Create Hadists

DROP INDEX IF EXISTS idx_hadist_bookmarks_hadist;
DROP INDEX IF EXISTS idx_hadist_bookmarks_user;
DROP TABLE IF EXISTS hadist_bookmarks;
DROP INDEX IF EXISTS idx_hadist_likes_hadist;
DROP INDEX IF EXISTS idx_hadist_likes_user;
DROP TABLE IF EXISTS hadist_likes;
DROP INDEX IF EXISTS idx_hadists_shahih;
DROP INDEX IF EXISTS idx_hadists_kitab;
DROP INDEX IF EXISTS idx_hadists_category;
DROP INDEX IF EXISTS idx_hadists_uuid;
DROP TABLE IF EXISTS hadists;
