-- Down Migration: Create Users Table

DROP PROCEDURE IF EXISTS sp_revoke_user_tokens_v2;
DROP FUNCTION IF EXISTS sp_check_user_lock_status;
DROP FUNCTION IF EXISTS sp_handle_user_login_failure;
DROP INDEX IF EXISTS idx_user_refresh_tokens_user_id;
DROP TABLE IF EXISTS user_refresh_tokens;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP TABLE IF EXISTS users;
