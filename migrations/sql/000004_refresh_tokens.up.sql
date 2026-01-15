-- Up Migration: Refresh Tokens SP

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
