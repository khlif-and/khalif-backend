-- Stored Procedure: sp_record_listening
-- Records user listening and increments count (max 1 per user per audio per day)

CREATE OR REPLACE FUNCTION sp_record_listening(
    p_user_id BIGINT,
    p_audio_id BIGINT
) RETURNS TABLE (
    already_listened BOOLEAN,
    new_count BIGINT
) AS $$
DECLARE
    v_exists BOOLEAN;
    v_new_count BIGINT;
BEGIN
    -- Check if user already listened to this audio today
    SELECT EXISTS(
        SELECT 1 FROM listening_histories
        WHERE user_id = p_user_id 
        AND audio_id = p_audio_id
        AND listened_at >= CURRENT_DATE
    ) INTO v_exists;

    IF v_exists THEN
        -- Already listened today, just return current count
        SELECT listening_count INTO v_new_count FROM audios WHERE id = p_audio_id;
        RETURN QUERY SELECT TRUE, v_new_count;
    ELSE
        -- Record new listening
        INSERT INTO listening_histories (user_id, audio_id, listened_at)
        VALUES (p_user_id, p_audio_id, NOW());
        
        -- Increment listening count
        UPDATE audios 
        SET listening_count = listening_count + 1 
        WHERE id = p_audio_id
        RETURNING listening_count INTO v_new_count;
        
        RETURN QUERY SELECT FALSE, v_new_count;
    END IF;
END;
$$ LANGUAGE plpgsql;
