-- Migration: Add Google ID to Users table

ALTER TABLE users ADD COLUMN google_id VARCHAR(100);
CREATE UNIQUE INDEX idx_users_google_id ON users(google_id);
