-- Up Migration: Create Users Table

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    profile_picture TEXT,
    failed_login_attempts INT DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Create user_refresh_tokens table
CREATE TABLE IF NOT EXISTS user_refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    user_agent VARCHAR(255),
    ip_address VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_refresh_token_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_refresh_tokens_user_id ON user_refresh_tokens(user_id);

-- Create stored procedure for user login failures
CREATE OR REPLACE FUNCTION sp_handle_user_login_failure(
    p_email VARCHAR
) RETURNS TABLE (
    is_locked BOOLEAN,
    remaining_lock_time INTERVAL
) AS $$
DECLARE
    v_failed_attempts INT;
    v_locked_until TIMESTAMP WITH TIME ZONE;
    v_max_attempts INT := 5;
    v_lock_duration INTERVAL := '30 minutes';
BEGIN
    SELECT failed_login_attempts, locked_until
    INTO v_failed_attempts, v_locked_until
    FROM users
    WHERE email = p_email AND deleted_at IS NULL;

    IF NOT FOUND THEN
        RETURN;
    END IF;

    IF v_locked_until IS NOT NULL AND v_locked_until > NOW() THEN
        RETURN QUERY SELECT TRUE, v_locked_until - NOW();
        RETURN;
    END IF;

    UPDATE users
    SET failed_login_attempts = failed_login_attempts + 1,
        locked_until = CASE
            WHEN failed_login_attempts + 1 >= v_max_attempts THEN NOW() + v_lock_duration
            ELSE NULL
        END
    WHERE email = p_email;

    SELECT locked_until IS NOT NULL AND locked_until > NOW(),
           CASE WHEN locked_until > NOW() THEN locked_until - NOW() ELSE NULL END
    INTO is_locked, remaining_lock_time
    FROM users
    WHERE email = p_email;
    
    RETURN QUERY SELECT is_locked, remaining_lock_time;
END;
$$ LANGUAGE plpgsql;

-- Create stored procedure for checking user lock status
CREATE OR REPLACE FUNCTION sp_check_user_lock_status(
    p_email VARCHAR
) RETURNS TABLE (
    is_locked BOOLEAN,
    remaining_lock_time INTERVAL
) AS $$
DECLARE
    v_locked_until TIMESTAMP WITH TIME ZONE;
BEGIN
    SELECT locked_until INTO v_locked_until
    FROM users
    WHERE email = p_email AND deleted_at IS NULL;

    IF v_locked_until IS NOT NULL AND v_locked_until > NOW() THEN
        RETURN QUERY SELECT TRUE, v_locked_until - NOW();
    ELSE
        RETURN QUERY SELECT FALSE, NULL::INTERVAL;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create stored procedure to revoke user tokens
CREATE OR REPLACE PROCEDURE sp_revoke_user_tokens_v2(
    p_user_id INT,
    p_jti VARCHAR,
    p_is_rotation BOOLEAN DEFAULT FALSE
) AS $$
BEGIN
    IF p_jti IS NOT NULL AND p_jti != '' THEN
        UPDATE user_refresh_tokens
        SET is_revoked = TRUE
        WHERE user_id = p_user_id AND token_hash = p_jti;
    END IF;
    
    IF NOT p_is_rotation THEN
        UPDATE user_refresh_tokens
        SET is_revoked = TRUE
        WHERE user_id = p_user_id;
    END IF;
END;
$$ LANGUAGE plpgsql;
