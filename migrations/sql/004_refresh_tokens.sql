-- Up Migration
-- Table Schema managed by GORM AutoMigrate (Hybrid Approach)

-- Stored Procedure to Revoke All Tokens for a User (Security)
CREATE OR REPLACE PROCEDURE sp_revoke_user_tokens(
    p_admin_id INT
)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE refresh_tokens
    SET is_revoked = TRUE
    WHERE admin_id = p_admin_id AND is_revoked = FALSE;
END;
$$;

-- Down Migration
-- DROP TABLE IF EXISTS refresh_tokens;
-- DROP PROCEDURE IF EXISTS sp_revoke_user_tokens;
